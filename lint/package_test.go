package lint_test

import (
	"go/ast"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/mgechev/revive/lint"
)

// varTypesRule reports the type of the value of every package-level variable.
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
	// Make the result independent of the environment of the test.
	t.Setenv("GOFLAGS", "")
	t.Setenv("GOWORK", "off")

	rootDir := t.TempDir()
	files := map[string]string{
		"go.mod":  "module example.com/m\n\ngo 1.21\n",
		"b/b.go":  "package b\n\ntype T struct{}\n\nfunc New() T { return T{} }\n",
		"a/a.go":  "package a\n\nimport (\n\t\"strings\"\n\n\t\"example.com/m/b\"\n)\n\nvar (\n\tX = b.New()\n\tY = strings.ToUpper(\"\")\n)\n",
		"a/a2.go": "package a\n\nimport \"example.com/m/b\"\n\nvar Z = []b.T{}\n",
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
	pkgA := []string{filepath.Join(rootDir, "a", "a.go"), filepath.Join(rootDir, "a", "a2.go")}
	pkgB := []string{filepath.Join(rootDir, "b", "b.go")}
	failures, err := l.Lint([][]string{pkgA, pkgB}, []lint.Rule{&varTypesRule{}}, lint.Config{})
	if err != nil {
		t.Fatal("unexpected error from linting:", err)
	}

	var got []string
	for failure := range failures {
		got = append(got, failure.Failure)
	}
	slices.Sort(got)

	want := []string{"[]example.com/m/b.T", "example.com/m/b.T", "string"}
	if !slices.Equal(got, want) {
		t.Errorf("got types %q, want %q", got, want)
	}
}
