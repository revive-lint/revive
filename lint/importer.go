package lint

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// moduleExports locates the compiled export data of the packages imported by the linted code of a module.
//
// It asks the go command, so it doesn't depend on GOROOT/GOPATH being set or
// on the revive binary being built without -trimpath (see issue #1277).
// Export data is resolved lazily, in batches: the first missing import path triggers
// a single go list invocation that resolves the dependencies of all the linted packages
// of the module not loaded yet. Results are cached for the life of the process,
// as the standard library importer does.
type moduleExports struct {
	dir string // directory where the go command runs

	mu      sync.Mutex
	known   map[string]bool         // linted package directories
	pending []string                // linted package directories whose dependencies aren't loaded yet
	exports map[string]exportResult // import path -> export data file
}

type exportResult struct {
	file string
	err  error
}

// modulesExports caches *moduleExports by module root directory.
var modulesExports sync.Map

// newExportResolvers returns, for each linted package, the export data resolver of its module.
func newExportResolvers(packages [][]string) ([]*moduleExports, error) {
	resolvers := make([]*moduleExports, len(packages))
	byDir := map[string]*moduleExports{}
	for n, files := range packages {
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
			v, _ := modulesExports.LoadOrStore(root, &moduleExports{dir: root})
			m = v.(*moduleExports)
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

	return os.Open(file) //nolint:gosec // ignore G304: the file comes from the go command
}

func (m *moduleExports) exportFile(path string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if res, ok := m.exports[path]; ok {
		return res.file, res.err
	}

	if len(m.pending) > 0 {
		m.loadPending()
		if res, ok := m.exports[path]; ok {
			return res.file, res.err
		}
	}

	// This happens, for example, when the linted files don't form a valid package.
	out, err := m.goList("list", "-export", "-f", "{{.Export}}", "--", path)
	file := string(bytes.TrimSpace(out))
	if err == nil && file == "" {
		err = fmt.Errorf("no export data for %q", path)
	}
	m.exports[path] = exportResult{file: file, err: err}

	return file, err
}

// listedPackage holds the fields of the go list JSON output used by moduleExports.
type listedPackage struct {
	ImportPath string            `json:"ImportPath"`
	Name       string            `json:"Name"`
	Export     string            `json:"Export"`
	ImportMap  map[string]string `json:"ImportMap"`
}

// loadPending resolves the export data of all dependencies (including test dependencies) of the pending packages.
//
// It runs in two steps to avoid compiling the test variants of the linted packages,
// which can't be imported anyway: the first step lists the dependencies without
// building anything, the second one builds the export data of these dependencies.
func (m *moduleExports) loadPending() {
	args := []string{"list", "-e", "-deps", "-test", "-json=ImportPath,Name,ImportMap"}
	for _, dir := range m.pending {
		args = append(args, m.relativeDir(dir))
	}
	m.pending = nil

	deps, err := m.listPackages(args...)
	if err != nil {
		return
	}

	importMap := map[string]string{}
	args = []string{"list", "-e", "-export", "-json=ImportPath,Export", "--"}
	for _, pkg := range deps {
		maps.Copy(importMap, pkg.ImportMap)

		// Skip test variants like "foo [foo.test]" and test mains like "foo.test":
		// they can't be referenced by an import path.
		isTestMain := pkg.Name == "main" && strings.HasSuffix(pkg.ImportPath, ".test")
		if strings.Contains(pkg.ImportPath, " ") || isTestMain {
			continue
		}
		if _, ok := m.exports[pkg.ImportPath]; !ok {
			args = append(args, pkg.ImportPath)
		}
	}

	exports, err := m.listPackages(args...)
	if err != nil {
		return
	}
	for _, pkg := range exports {
		if pkg.Export != "" {
			m.exports[pkg.ImportPath] = exportResult{file: pkg.Export}
		}
	}

	// Map the import paths as written in the source code to the actual packages (e.g. vendored packages).
	for from, to := range importMap {
		if res, ok := m.exports[to]; ok {
			if _, exists := m.exports[from]; !exists {
				m.exports[from] = res
			}
		}
	}
}

func (m *moduleExports) listPackages(args ...string) ([]listedPackage, error) {
	out, err := m.goList(args...)
	if err != nil {
		return nil, err
	}

	var pkgs []listedPackage
	dec := json.NewDecoder(bytes.NewReader(out))
	for {
		var pkg listedPackage
		if err := dec.Decode(&pkg); err != nil {
			if errors.Is(err, io.EOF) {
				return pkgs, nil
			}
			return pkgs, err
		}
		pkgs = append(pkgs, pkg)
	}
}

func (m *moduleExports) relativeDir(dir string) string {
	rel, err := filepath.Rel(m.dir, dir)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return dir
	}
	return "." + string(filepath.Separator) + rel
}

func (m *moduleExports) goList(args ...string) ([]byte, error) {
	cmd := exec.Command("go", args...) //nolint:gosec // ignore G204: the arguments are controlled by revive
	cmd.Dir = m.dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if _, ok := errors.AsType[*exec.ExitError](err); ok {
		return out, fmt.Errorf("go %s: %w: %s", strings.Join(args, " "), err, bytes.TrimSpace(stderr.Bytes()))
	}
	return out, err
}
