package lint_test

import (
	"go/ast"
	"go/types"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"testing"

	"github.com/mgechev/revive/lint"
)

// varTypesRule reports the types of variable and constant initializers.
type varTypesRule struct{}

func (*varTypesRule) Name() string { return "var-types" }

func (*varTypesRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if err := file.Pkg.TypeCheck(); err != nil {
		return []lint.Failure{lint.NewInternalFailure(err.Error())}
	}

	var failures []lint.Failure
	ast.Inspect(file.AST, func(n ast.Node) bool {
		spec, ok := n.(*ast.ValueSpec)
		if !ok {
			return true
		}
		for _, value := range spec.Values {
			typ := "<nil>"
			if t := file.Pkg.TypeOf(value); t != nil {
				typ = t.String()
			}
			failures = append(failures, lint.Failure{Confidence: 1, Node: value, Failure: typ})
		}
		return false
	})
	return failures
}

func TestPackage_TypeCheckResolvesModuleImports(t *testing.T) {
	// Keep module resolution independent of the caller's workspace and build flags.
	t.Setenv("GOFLAGS", "")
	t.Setenv("GOWORK", "off")

	rootDir := t.TempDir()
	files := map[string]string{
		"go.mod": "module example.com/m\n\ngo 1.21\n",
		"b/b.go": `package b

type T struct{}

func New() T { return T{} }
`,
		"a/a.go": `package a

import (
	"strings"

	"example.com/m/b"
)

var (
	X = b.New()
	Y = strings.ToUpper("")
)
`,
		"a/a2.go": `package a

import "example.com/m/b"

var Z = []b.T{}
`,
		"internaldep/internaldep.go": `package internaldep

type T struct{}
`,
		"externaldep/externaldep.go": `package externaldep

type T struct{}
`,
		"a/a_test.go": `package a

import "example.com/m/internaldep"

var Internal = internaldep.T{}
`,
		"a/a_external_test.go": `package a_test

import (
	"example.com/m/a"
	"example.com/m/externaldep"
)

var (
	External = externaldep.T{}
	Production = a.X
)
`,
	}
	writeFiles(t, rootDir, files)

	l := lint.New(os.ReadFile, 0)
	pkgA := []string{
		filepath.Join(rootDir, "a", "a.go"),
		filepath.Join(rootDir, "a", "a2.go"),
		filepath.Join(rootDir, "a", "a_test.go"),
	}
	// External tests form a separate package even though they share a directory.
	pkgATest := []string{filepath.Join(rootDir, "a", "a_external_test.go")}
	pkgB := []string{filepath.Join(rootDir, "b", "b.go")}
	failures, err := l.Lint([][]string{pkgA, pkgATest, pkgB}, []lint.Rule{&varTypesRule{}}, lint.Config{})
	if err != nil {
		t.Fatal("unexpected error from linting:", err)
	}

	var got []string
	for failure := range failures {
		got = append(got, failure.Failure)
	}
	slices.Sort(got)

	want := []string{
		"[]example.com/m/b.T",
		"example.com/m/b.T",
		"example.com/m/b.T",
		"example.com/m/externaldep.T",
		"example.com/m/internaldep.T",
		"string",
	}
	if !slices.Equal(got, want) {
		t.Errorf("got types %q, want %q", got, want)
	}
}

// importedPkgsRule records the package each linted package imports as importPath.
type importedPkgsRule struct {
	importPath string

	mu       sync.Mutex
	imported map[*types.Package]bool
}

func (*importedPkgsRule) Name() string { return "imported-pkgs" }

func (r *importedPkgsRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if err := file.Pkg.TypeCheck(); err != nil {
		return []lint.Failure{lint.NewInternalFailure(err.Error())}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	for _, imp := range file.Pkg.TypesPkg().Imports() {
		if imp.Path() == r.importPath {
			r.imported[imp] = true
		}
	}
	return nil
}

func TestPackage_TypeCheckSharesImportedPackages(t *testing.T) {
	// Keep module resolution independent of the caller's workspace and build flags.
	t.Setenv("GOFLAGS", "")
	t.Setenv("GOWORK", "off")

	rootDir := t.TempDir()
	writeFiles(t, rootDir, map[string]string{
		"go.mod": "module example.com/m\n\ngo 1.21\n",
		"b/b.go": "package b\n\ntype T struct{}\n",
		"a/a.go": "package a\n\nimport \"example.com/m/b\"\n\nvar X b.T\n",
		"c/c.go": "package c\n\nimport \"example.com/m/b\"\n\nvar Y b.T\n",
	})

	rule := &importedPkgsRule{importPath: "example.com/m/b", imported: map[*types.Package]bool{}}
	l := lint.New(os.ReadFile, 0)
	failures, err := l.Lint([][]string{
		{filepath.Join(rootDir, "a", "a.go")},
		{filepath.Join(rootDir, "c", "c.go")},
	}, []lint.Rule{rule}, lint.Config{})
	if err != nil {
		t.Fatal("unexpected error from linting:", err)
	}
	for failure := range failures {
		t.Errorf("unexpected failure: %s", failure.Failure)
	}

	// Packages of the same module decode a common dependency once and share it.
	if got := len(rule.imported); got != 1 {
		t.Errorf("got %d distinct imported packages for %q, want 1", got, rule.importPath)
	}
}

// writeFiles writes files, keyed by slash-separated path relative to rootDir.
func writeFiles(t *testing.T, rootDir string, files map[string]string) {
	t.Helper()

	for name, content := range files {
		path := filepath.Join(rootDir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}
