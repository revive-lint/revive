package astutils

import (
	"go/ast"
	"strconv"
)

// exitFuncChecker is a function type that checks whether a function call is an exit function.
type exitFuncChecker func(args []ast.Expr, aliases ExitPkgAliases) bool

var alwaysTrue exitFuncChecker = func([]ast.Expr, ExitPkgAliases) bool { return true }

// exitFunctions is a map of std packages and functions that are considered as exit functions.
var exitFunctions = map[string]map[string]exitFuncChecker{
	"os":      {"Exit": alwaysTrue},
	"syscall": {"Exit": alwaysTrue},
	"log": {
		"Fatal":   alwaysTrue,
		"Fatalf":  alwaysTrue,
		"Fatalln": alwaysTrue,
		"Panic":   alwaysTrue,
		"Panicf":  alwaysTrue,
		"Panicln": alwaysTrue,
	},
	"flag": {
		"Parse": func([]ast.Expr, ExitPkgAliases) bool { return true },
		"NewFlagSet": func(args []ast.Expr, aliases ExitPkgAliases) bool {
			if len(args) != 2 {
				return false
			}
			sel, ok := args[1].(*ast.SelectorExpr)
			if !ok {
				return false
			}
			id, ok := sel.X.(*ast.Ident)
			return ok && aliases.PkgName(id.Name) == "flag" && sel.Sel.Name == "ExitOnError"
		},
	},
}

// ExitPkgAliases maps the names of renamed imports of the std packages that declare exit functions
// to the names of these packages, e.g. "sys" to "os" for import sys "os".
type ExitPkgAliases map[string]string

// NewExitPkgAliases returns the renamed imports of the std packages that declare exit functions.
func NewExitPkgAliases(imports []*ast.ImportSpec) ExitPkgAliases {
	aliases := ExitPkgAliases{}
	for _, imp := range imports {
		if imp.Name == nil || imp.Name.Name == "_" || imp.Name.Name == "." {
			continue
		}

		path, err := strconv.Unquote(imp.Path.Value)
		if err != nil {
			continue
		}

		if _, ok := exitFunctions[path]; ok {
			aliases[imp.Name.Name] = path
		}
	}

	return aliases
}

// PkgName returns the name of the std package that name refers to if name is a renamed import of it, or name otherwise.
func (a ExitPkgAliases) PkgName(name string) string {
	if pkgName, ok := a[name]; ok {
		return pkgName
	}

	return name
}

// IsCallToExitFunction checks if the function call is a call to an exit function.
// The package name is the one used in the call; aliases resolves renamed imports and may be nil.
func IsCallToExitFunction(pkgName, functionName string, callArgs []ast.Expr, aliases ExitPkgAliases) bool {
	m, ok := exitFunctions[aliases.PkgName(pkgName)]
	if !ok {
		return false
	}

	check, ok := m[functionName]
	if !ok {
		return false
	}

	return check(callArgs, aliases)
}
