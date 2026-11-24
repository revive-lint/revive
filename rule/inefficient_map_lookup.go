package rule

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/mgechev/revive/internal/astutils"
	"github.com/mgechev/revive/lint"
)

// InefficientMapLookupRule spots potential inefficient map lookups.
type InefficientMapLookupRule struct{}

// Apply applies the rule to given file.
func (*InefficientMapLookupRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintInefficientMapLookup{
		file:      file,
		onFailure: onFailure,
	}

	if err := file.Pkg.TypeCheck(); err != nil {
		return []lint.Failure{
			lint.NewInternalFailure(fmt.Sprintf("Unable to type check file %q: %v", file.Name, err)),
		}
	}

	// Iterate over declarations looking for function declarations
	for _, decl := range file.AST.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue // not a function
		}

		if fn.Body == nil {
			continue // external (no-Go) function
		}

		// Analyze the function body
		ast.Walk(w, fn.Body)
	}

	return failures
}

// Name returns the rule name.
func (*InefficientMapLookupRule) Name() string {
	return "inefficient-map-lookup"
}

type lintInefficientMapLookup struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

func (w *lintInefficientMapLookup) Visit(node ast.Node) ast.Visitor {
	// Only interested in blocks of statements
	block, ok := node.(*ast.BlockStmt)
	if !ok {
		return w // not a block of statements
	}

	w.analyzeBlock(block)

	return w
}

// analyzeBlock searches AST subtrees with the following form
//
//	for <key> := range <map> {
//		if <key> == <something> {
//		   ...
//		}
func (w *lintInefficientMapLookup) analyzeBlock(b *ast.BlockStmt) {
	for _, stmt := range b.List {
		rangeOverMap, key := w.rangeOverMapKey(stmt)
		if rangeOverMap == nil {
			continue
		}

		// Here we have identified a range over the keys of a map
		// Let's check if the range body is
		// { if <key> == <something> { ... } }
		// or
		// { if <key> != <something> { continue } ... }
		if !isKeyLookup(key.Name, rangeOverMap.Body) {
			continue
		}

		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       rangeOverMap,
			Category:   lint.FailureCategoryStyle,
			Failure:    "inefficient lookup of map key",
		})
	}
}

func isKeyLookup(keyName string, blockStmt *ast.BlockStmt) bool {
	blockLen := len(blockStmt.List)
	if blockLen == 0 {
		return false // empty
	}

	firstStmt := blockStmt.List[0]
	ifStmt, ok := firstStmt.(*ast.IfStmt)
	if !ok {
		return false // the first statement of the body is not an if
	}

	binExp, ok := ifStmt.Cond.(*ast.BinaryExpr)
	if !ok {
		return false // the if condition is not a binary expression
	}

	if !astutils.IsIdent(binExp.X, keyName) {
		return false // the if condition is not <key> <bin-op> <LHS>
	}

	switch binExp.Op {
	case token.EQL:
		// if key == ... should be the single statement in the block
		return blockLen == 1

	case token.NEQ:
		// if key != ...
		ifBodyStmts := ifStmt.Body.List
		if len(ifBodyStmts) < 1 {
			return false // if key != ... { /* empty */ }
		}

		branchStmt, ok := ifBodyStmts[0].(*ast.BranchStmt)
		if !ok || branchStmt.Tok != token.CONTINUE {
			return false // if key != ... { <not a continue> }
		}

		return true
	}

	return false
}

// rangeOverMapKey returns the range statement and its key identifier when stmt
// ranges over the keys of a map, and nil otherwise.
func (w *lintInefficientMapLookup) rangeOverMapKey(stmt ast.Stmt) (*ast.RangeStmt, *ast.Ident) {
	rangeStmt, ok := stmt.(*ast.RangeStmt)
	if !ok {
		return nil, nil // not a range
	}

	// Check if we range on the key, and that the key is a plain identifier:
	// an assigning range like "for s.field = range aMap" has a non-identifier key.
	key, ok := rangeStmt.Key.(*ast.Ident)
	if !ok {
		return nil, nil
	}

	// Check if we range only on key
	// for key := range ...
	// for key, _ := range ...
	hasValueVariable := rangeStmt.Value != nil && !astutils.IsIdent(rangeStmt.Value, "_")
	if hasValueVariable {
		return nil, nil // range over both key and value
	}

	// Check if we range over a map
	t := w.file.Pkg.TypeOf(rangeStmt.X)
	if t == nil || !isMapType(t) {
		return nil, nil
	}

	return rangeStmt, key
}

// isMapType reports whether t is a map type, including named and alias map types
// and type parameters whose constraint only allows maps.
func isMapType(t types.Type) bool {
	switch u := t.Underlying().(type) {
	case *types.Map:
		return true
	case *types.Interface:
		// The underlying type of a type parameter is its constraint interface.
		for embedded := range u.EmbeddedTypes() {
			if isMapType(embedded) {
				return true
			}
		}
	case *types.Union:
		for term := range u.Terms() {
			if !isMapType(term.Type()) {
				return false
			}
		}
		return u.Len() > 0
	}
	return false
}
