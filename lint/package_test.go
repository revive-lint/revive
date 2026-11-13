package lint_test

import (
	"go/ast"
	"os"
	"path/filepath"
	"slices"
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
	for name, content := range files {
		path := filepath.Join(rootDir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}

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

func TestPackage_TypeCheckRefreshesModuleImports(t *testing.T) {
	t.Setenv("GOFLAGS", "")
	t.Setenv("GOWORK", "off")

	rootDir := t.TempDir()
	files := map[string]string{
		"go.mod": "module example.com/m\n\ngo 1.21\n",
		"a/a.go": "package a\nimport \"example.com/m/b\"\nvar X = b.Value\n",
		"b/b.go": "package b\nvar Value int\n",
	}
	for name, content := range files {
		path := filepath.Join(rootDir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	l := lint.New(os.ReadFile, 0)
	pkgs := [][]string{{filepath.Join(rootDir, "a", "a.go")}}
	rules := []lint.Rule{&varTypesRule{}}
	for _, want := range []string{"int", "string"} {
		t.Run(want, func(t *testing.T) {
			// Reuse the linter after changing an imported package on disk.
			source := "package b\nvar Value " + want + "\n"
			if err := os.WriteFile(filepath.Join(rootDir, "b", "b.go"), []byte(source), 0o600); err != nil {
				t.Fatal(err)
			}
			failures, err := l.Lint(pkgs, rules, lint.Config{})
			if err != nil {
				t.Fatal("unexpected error from linting:", err)
			}
			var got []string
			for failure := range failures {
				got = append(got, failure.Failure)
			}
			if !slices.Equal(got, []string{want}) {
				t.Errorf("got types %q, want %q", got, []string{want})
			}
		})
	}
}
