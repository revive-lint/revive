package lint

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/tools/go/packages"
)

// moduleExports locates compiled export data for imports of the linted packages.
//
// It uses [packages.Load] in the target module's build environment, including when
// revive is built with -trimpath. Dependencies are loaded lazily in batches and
// cached for the lifetime of the process.
type moduleExports struct {
	dir string

	mu      sync.Mutex
	known   map[string]bool         // Registered package directories.
	pending []string                // Package directories whose dependencies have not been loaded.
	exports map[string]exportResult // Export data files and errors indexed by import path.
}

type exportResult struct {
	file string
	err  error
}

// moduleExportsCache caches export resolvers by module root directory.
var moduleExportsCache sync.Map

// newExportResolvers returns the module export resolver for each linted package.
func newExportResolvers(pkgs [][]string) ([]*moduleExports, error) {
	resolvers := make([]*moduleExports, len(pkgs))
	byDir := map[string]*moduleExports{}
	for n, files := range pkgs {
		if len(files) == 0 {
			continue
		}

		dir, err := filepath.Abs(filepath.Dir(files[0]))
		if err != nil {
			return nil, err
		}

		m, ok := byDir[dir]
		if !ok {
			root := dir
			if modFile, err := retrieveModFile(dir); err == nil {
				root = filepath.Dir(modFile)
			}
			v, _ := moduleExportsCache.LoadOrStore(root, &moduleExports{dir: root})
			m, ok = v.(*moduleExports)
			if !ok {
				return nil, fmt.Errorf("invalid export resolver for module %q: got %T", root, v)
			}
			m.addPackage(dir)
			byDir[dir] = m
		}
		resolvers[n] = m
	}
	return resolvers, nil
}

func (m *moduleExports) addPackage(dir string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.known[dir] {
		return
	}
	if m.known == nil {
		m.known = map[string]bool{}
		m.exports = map[string]exportResult{}
	}
	m.known[dir] = true
	m.pending = append(m.pending, dir)
}

// lookup implements [go/importer.Lookup].
func (m *moduleExports) lookup(path string) (io.ReadCloser, error) {
	file, err := m.exportFile(path)
	if err != nil {
		return nil, err
	}
	return os.OpenInRoot(filepath.Dir(file), filepath.Base(file))
}

func (m *moduleExports) exportFile(path string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if res, ok := m.exports[path]; ok {
		return res.file, res.err
	}
	var batchErr error
	if len(m.pending) > 0 {
		batchErr = m.loadPending()
		if res, ok := m.exports[path]; ok {
			return res.file, res.err
		}
	}

	// Individual file sets may not form a buildable package. Resolve their
	// imports directly when dependency discovery did not find them.
	pkgs, err := packages.Load(&packages.Config{
		Dir:  m.dir,
		Mode: packages.NeedName | packages.NeedExportFile,
	}, path)
	var file string
	for _, pkg := range pkgs {
		if pkg.PkgPath != path {
			continue
		}
		file = pkg.ExportFile
		if file != "" {
			break
		}
		for _, pkgErr := range pkg.Errors {
			err = errors.Join(err, pkgErr)
		}
	}
	if file == "" {
		err = errors.Join(batchErr, err)
		if err == nil {
			err = fmt.Errorf("no export data for %q: package could not be loaded", path)
		} else {
			err = fmt.Errorf("no export data for %q: %w", path, err)
		}
	}
	m.exports[path] = exportResult{file: file, err: err}
	return file, err
}

// loadPending resolves export data for pending packages and their dependencies.
// It discovers test imports separately to avoid compiling test variants and generated test mains.
func (m *moduleExports) loadPending() error {
	patterns := make([]string, 0, len(m.pending))
	for _, dir := range m.pending {
		patterns = append(patterns, m.relativeDir(dir))
	}
	m.pending = nil
	pkgs, err := packages.Load(&packages.Config{
		Dir:   m.dir,
		Mode:  packages.NeedName | packages.NeedImports | packages.NeedDeps,
		Tests: true,
	}, patterns...)
	if err != nil {
		return fmt.Errorf("discover package dependencies: %w", err)
	}

	var discoveryErr error
	importMap := map[string]string{}
	paths := map[string]bool{}
	patterns = nil
	packages.Visit(pkgs, func(pkg *packages.Package) bool {
		for _, pkgErr := range pkg.Errors {
			discoveryErr = errors.Join(discoveryErr, pkgErr)
		}
		for path, imported := range pkg.Imports {
			importMap[path] = imported.PkgPath
		}
		isTestMain := pkg.Name == "main" && strings.HasSuffix(pkg.PkgPath, ".test")
		if pkg.ID != pkg.PkgPath || isTestMain {
			return true
		}
		if _, ok := m.exports[pkg.PkgPath]; !ok && !paths[pkg.PkgPath] {
			paths[pkg.PkgPath] = true
			patterns = append(patterns, pkg.PkgPath)
		}
		return true
	}, nil)
	if discoveryErr != nil {
		discoveryErr = fmt.Errorf("discover package dependencies: %w", discoveryErr)
	}
	if len(patterns) == 0 {
		return discoveryErr
	}

	pkgs, err = packages.Load(&packages.Config{
		Dir:  m.dir,
		Mode: packages.NeedName | packages.NeedExportFile,
	}, patterns...)
	if err != nil {
		return errors.Join(discoveryErr, fmt.Errorf("load package export data: %w", err))
	}
	var exportErr error
	for _, pkg := range pkgs {
		if pkg.ExportFile != "" {
			m.exports[pkg.PkgPath] = exportResult{file: pkg.ExportFile}
		}
		for _, pkgErr := range pkg.Errors {
			exportErr = errors.Join(exportErr, pkgErr)
		}
	}

	// Vendored packages in GOPATH mode can have a package path different from
	// the import string used in the source.
	for from, to := range importMap {
		if res, ok := m.exports[to]; ok {
			if _, exists := m.exports[from]; !exists {
				m.exports[from] = res
			}
		}
	}
	if exportErr != nil {
		exportErr = fmt.Errorf("load package export data: %w", exportErr)
	}
	return errors.Join(discoveryErr, exportErr)
}

func (m *moduleExports) relativeDir(dir string) string {
	rel, err := filepath.Rel(m.dir, dir)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return dir
	}
	return "." + string(filepath.Separator) + rel
}
