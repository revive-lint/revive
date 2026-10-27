package lint

import (
	"errors"
	"go/ast"
	"go/importer"
	"go/token"
	"go/types"
	"log/slog"
	"slices"
	"sync"

	goversion "github.com/hashicorp/go-version"
	"golang.org/x/sync/errgroup"

	"github.com/mgechev/revive/internal/astutils"
	"github.com/mgechev/revive/internal/timing"
	"github.com/mgechev/revive/internal/typeparams"
)

// Package represents a package in the project.
type Package struct {
	fset *token.FileSet

	mu        sync.RWMutex
	files     map[string]*File
	goVersion *goversion.Version
	typesPkg  *types.Package
	typesInfo *types.Info
	// sortable is the set of types in the package that implement sort.Interface.
	sortable map[string]bool
	// main is whether this is a "main" package.
	main int

	logger *slog.Logger
}

var (
	trueValue  = 1
	falseValue = 2

	// Go111 is a constant representing the Go version 1.11.
	Go111 = goversion.Must(goversion.NewVersion("1.11"))
	// Go115 is a constant representing the Go version 1.15.
	Go115 = goversion.Must(goversion.NewVersion("1.15"))
	// Go118 is a constant representing the Go version 1.18.
	Go118 = goversion.Must(goversion.NewVersion("1.18"))
	// Go121 is a constant representing the Go version 1.21.
	Go121 = goversion.Must(goversion.NewVersion("1.21"))
	// Go122 is a constant representing the Go version 1.22.
	Go122 = goversion.Must(goversion.NewVersion("1.22"))
	// Go124 is a constant representing the Go version 1.24.
	Go124 = goversion.Must(goversion.NewVersion("1.24"))
	// Go125 is a constant representing the Go version 1.25.
	Go125 = goversion.Must(goversion.NewVersion("1.25"))
	// Go126 is a constant representing the Go version 1.26.
	Go126 = goversion.Must(goversion.NewVersion("1.26"))
)

// Files return package's files.
func (p *Package) Files() map[string]*File {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.files
}

// IsMain returns if that's the main package.
func (p *Package) IsMain() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	switch p.main {
	case trueValue:
		return true
	case falseValue:
		return false
	}
	for _, f := range p.files {
		if f.isMain() {
			p.main = trueValue
			return true
		}
	}
	p.main = falseValue
	return false
}

// TypesPkg yields information on this package.
func (p *Package) TypesPkg() *types.Package {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.typesPkg
}

// TypesInfo yields type information of this package identifiers.
func (p *Package) TypesInfo() *types.Info {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.typesInfo
}

// Sortable yields a map of sortable types in this package.
func (p *Package) Sortable() map[string]bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.sortable
}

// TypeCheck performs type checking for given package.
func (p *Package) TypeCheck() error {
	// Deferred before the unlock to log once p.mu is released: rules waiting for it must not also wait for the log output.
	var logTiming func()
	defer func() {
		if logTiming != nil {
			logTiming()
		}
	}()

	p.mu.Lock()
	defer p.mu.Unlock()

	alreadyTypeChecked := p.typesInfo != nil || p.typesPkg != nil
	if alreadyTypeChecked {
		return nil
	}

	config := &types.Config{
		// By setting a no-op error reporter, the type checker does as much work as possible.
		Error:    func(error) {},
		Importer: importer.Default(),
	}
	info := &types.Info{
		Types:      map[ast.Expr]types.TypeAndValue{},
		Defs:       map[*ast.Ident]types.Object{},
		Uses:       map[*ast.Ident]types.Object{},
		Scopes:     map[ast.Node]*types.Scope{},
		Selections: map[*ast.SelectorExpr]*types.Selection{},
	}
	var anyFile *File
	var astFiles []*ast.File
	filenames := make([]string, 0, len(p.files))
	for name := range p.files {
		filenames = append(filenames, name)
	}
	slices.Sort(filenames)

	for _, name := range filenames {
		f := p.files[name]
		if anyFile == nil || (!f.IsTest() && anyFile.IsTest()) {
			anyFile = f
		}
		astFiles = append(astFiles, f.AST)
	}

	if anyFile == nil {
		// this is unlikely to happen, but technically guarantees anyFile to not be nil
		return errors.New("no ast.File found")
	}

	done := timing.Track(p.logger, "type-check")
	typesPkg, err := config.Check(anyFile.AST.Name.Name, p.fset, astFiles, info)
	if timing.Enabled(p.logger) {
		logTiming = func() { done("package", anyFile.AST.Name.Name, "files", len(astFiles), "error", err) }
	}

	// Remember the typechecking info, even if config.Check failed,
	// since we will get partial information.
	p.typesPkg = typesPkg
	p.typesInfo = info

	return err
}

// TypeOf returns the type of expression.
func (p *Package) TypeOf(expr ast.Expr) types.Type {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.typesInfo == nil {
		return nil
	}

	return p.typesInfo.TypeOf(expr)
}

type sortableMethodsFlags int

// flags for sortable interface methods.
const (
	bfLen sortableMethodsFlags = 1 << iota
	bfLess
	bfSwap
)

// scanSortable collects the sortable types of the package and returns their number.
func (p *Package) scanSortable() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	sortableFlags := map[string]sortableMethodsFlags{}
	for _, f := range p.files {
		for _, decl := range f.AST.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			isAMethodDeclaration := ok && fn.Recv != nil && len(fn.Recv.List) != 0
			if !isAMethodDeclaration {
				continue
			}

			recvType := typeparams.ReceiverType(fn)
			sortableFlags[recvType] |= getSortableMethodFlagForFunction(fn)
		}
	}

	p.sortable = make(map[string]bool, len(sortableFlags))
	for typ, ms := range sortableFlags {
		if ms == bfLen|bfLess|bfSwap {
			p.sortable[typ] = true
		}
	}
	return len(p.sortable)
}

func (p *Package) lint(rules []Rule, config Config, failures chan Failure) error {
	doneSortable := timing.Track(p.logger, "scan-sortable")
	sortable := p.scanSortable()
	if timing.Enabled(p.logger) {
		doneSortable("sortable", sortable)
	}

	var eg errgroup.Group
	for _, file := range p.Files() {
		eg.Go(func() error {
			return file.lint(rules, config, failures)
		})
	}

	return eg.Wait()
}

// IsAtLeastGoVersion returns true if the Go version for this package is v or higher, false otherwise.
func (p *Package) IsAtLeastGoVersion(v *goversion.Version) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.goVersion.GreaterThanOrEqual(v)
}

// GoVersion returns the Go version for this package.
func (p *Package) GoVersion() *goversion.Version {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.goVersion
}

func getSortableMethodFlagForFunction(fn *ast.FuncDecl) sortableMethodsFlags {
	switch {
	case astutils.FuncSignatureIs(fn, "Len", []string{}, []string{"int"}):
		return bfLen
	case astutils.FuncSignatureIs(fn, "Less", []string{"int", "int"}, []string{"bool"}):
		return bfLess
	case astutils.FuncSignatureIs(fn, "Swap", []string{"int", "int"}, []string{}):
		return bfSwap
	default:
		return 0
	}
}
