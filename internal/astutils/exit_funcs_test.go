package astutils_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"maps"
	"testing"

	"github.com/mgechev/revive/internal/astutils"
)

func TestIsCallToExitFunction(t *testing.T) {
	flagAliases := astutils.ExitPkgAliases{"fl": "flag"}
	tests := []struct {
		pkgName      string
		functionName string
		functionArgs []ast.Expr
		aliases      astutils.ExitPkgAliases
		expected     bool
	}{
		{"os", "Exit", nil, nil, true},
		{"syscall", "Exit", nil, nil, true},
		{"log", "Fatal", nil, nil, true},
		{"log", "Fatalf", nil, nil, true},
		{"log", "Fatalln", nil, nil, true},
		{"log", "Panic", nil, nil, true},
		{"log", "Panicf", nil, nil, true},
		{"flag", "Parse", nil, nil, true},
		{"flag", "NewFlagSet", []ast.Expr{
			nil,
			&ast.SelectorExpr{
				X:   &ast.Ident{Name: "flag"},
				Sel: &ast.Ident{Name: "ExitOnError"},
			},
		}, nil, true},
		{"log", "Print", nil, nil, false},
		{"fmt", "Errorf", nil, nil, false},
		{"flag", "NewFlagSet", []ast.Expr{
			nil,
			&ast.SelectorExpr{
				X:   &ast.Ident{Name: "flag"},
				Sel: &ast.Ident{Name: "ContinueOnError"},
			},
		}, nil, false},
		{"sys", "Exit", nil, astutils.ExitPkgAliases{"sys": "os"}, true},
		{"sys", "Exit", nil, nil, false},
		{"stdlog", "Fatal", nil, astutils.ExitPkgAliases{"stdlog": "log"}, true},
		{"stdlog", "Print", nil, astutils.ExitPkgAliases{"stdlog": "log"}, false},
		{"fl", "NewFlagSet", []ast.Expr{
			nil,
			&ast.SelectorExpr{
				X:   &ast.Ident{Name: "fl"},
				Sel: &ast.Ident{Name: "ExitOnError"},
			},
		}, flagAliases, true},
		{"fl", "NewFlagSet", []ast.Expr{
			nil,
			&ast.SelectorExpr{
				X:   &ast.Ident{Name: "fl"},
				Sel: &ast.Ident{Name: "ContinueOnError"},
			},
		}, flagAliases, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s.%s", tt.pkgName, tt.functionName), func(t *testing.T) {
			if got := astutils.IsCallToExitFunction(tt.pkgName, tt.functionName, tt.functionArgs, tt.aliases); got != tt.expected {
				t.Errorf("IsCallToExitFunction(%s, %s, %v, %v) = %v; want %v", tt.pkgName, tt.functionName, tt.functionArgs, tt.aliases, got, tt.expected)
			}
		})
	}
}

func TestNewExitPkgAliases(t *testing.T) {
	const src = `package p

import (
	"flag"
	_ "log"
	. "syscall"
	sys "os"
	log "example.com/logger"
)
`
	file, err := parser.ParseFile(token.NewFileSet(), "p.go", src, parser.ImportsOnly)
	if err != nil {
		t.Fatal(err)
	}

	got := astutils.NewExitPkgAliases(file.Imports)
	want := astutils.ExitPkgAliases{"sys": "os"}
	if !maps.Equal(got, want) {
		t.Errorf("NewExitPkgAliases() = %v; want %v", got, want)
	}
}
