package lint

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"sync"
	"testing"
)

func TestPackageTypeCheckWithEmptyGOROOT(t *testing.T) {
	t.Setenv("GOROOT", "")

	fset := token.NewFileSet()
	src := `package p
import "fmt"
func f() { fmt.Println("ok") }
`

	parsed, err := parser.ParseFile(fset, "p.go", src, parser.AllErrors)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	pkg := &Package{
		fset: fset,
		files: map[string]*File{
			"p.go": {
				Name: "p.go",
				AST:  parsed,
				Pkg:  nil,
			},
		},
	}
	pkg.files["p.go"].Pkg = pkg

	if err := pkg.TypeCheck(); err != nil {
		t.Fatalf("type check failed with empty GOROOT: %v", err)
	}

	var printCall *ast.CallExpr
	ast.Inspect(parsed, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if ok {
			printCall = call
			return false
		}
		return true
	})
	if printCall == nil {
		t.Fatal("expected to find fmt.Println call")
	}
	if got := pkg.TypeOf(printCall.Fun); got == nil {
		t.Fatal("expected type info for fmt.Println call")
	}
}

func TestEnsureSourceImporterGOROOTUsesGoEnvFallback(t *testing.T) {
	t.Setenv("GOROOT", "")

	prevCmd := sourceImporterGOROOTCmd
	prevResolved := sourceImporterGOROOT
	prevBuildGOROOT := build.Default.GOROOT
	t.Cleanup(func() {
		sourceImporterGOROOTCmd = prevCmd
		sourceImporterGOROOT = prevResolved
		sourceImporterGOROOTOnce = sync.Once{}
		build.Default.GOROOT = prevBuildGOROOT
	})

	sourceImporterGOROOTOnce = sync.Once{}
	sourceImporterGOROOT = ""
	sourceImporterGOROOTCmd = func() (string, error) {
		return "C:/go-fallback\n", nil
	}
	build.Default.GOROOT = ""

	ensureSourceImporterGOROOT("")

	if got := os.Getenv("GOROOT"); got != "C:/go-fallback" {
		t.Fatalf("expected GOROOT to be set from go env fallback, got %q", got)
	}
	if got := build.Default.GOROOT; got != "C:/go-fallback" {
		t.Fatalf("expected build.Default.GOROOT to be set from go env fallback, got %q", got)
	}
}
