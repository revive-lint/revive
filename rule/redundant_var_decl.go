package rule

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/mgechev/revive/lint"
)

// RedundantVarDeclRule warns on variable declarations made redundant by a later short variable declaration
// that redeclares the same variable in the same block.
type RedundantVarDeclRule struct {
	// onlyTypes restricts the rule to variables of these types. Empty means all types.
	onlyTypes map[string]struct{}
}

var _ lint.ConfigurableRule = (*RedundantVarDeclRule)(nil)

// Configure validates the rule configuration, and configures the rule accordingly.
//
// Configuration implements the [lint.ConfigurableRule] interface.
func (r *RedundantVarDeclRule) Configure(arguments lint.Arguments) error {
	r.onlyTypes = nil

	if len(arguments) == 0 {
		return nil
	}

	args, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf("invalid argument to the redundant-var-decl rule. Expecting a k,v map, got %T", arguments[0])
	}

	for k, v := range args {
		if !isRuleOption(k, "onlyTypes") {
			return fmt.Errorf("invalid argument to the redundant-var-decl rule: unknown option %q", k)
		}

		list, ok := v.([]any)
		if !ok {
			return fmt.Errorf("invalid argument to the redundant-var-decl rule, []string expected for %q. Got '%v' (%T)", k, v, v)
		}

		for _, item := range list {
			typeName, ok := item.(string)
			if !ok {
				return fmt.Errorf("invalid argument to the redundant-var-decl rule, string expected in %q. Got '%v' (%T)", k, item, item)
			}
			if r.onlyTypes == nil {
				r.onlyTypes = map[string]struct{}{}
			}
			r.onlyTypes[typeName] = struct{}{}
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *RedundantVarDeclRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck() //nolint:errcheck,gosec // G104: Errors unhandled
	info := file.Pkg.TypesInfo()
	if info == nil {
		return nil
	}

	w := lintRedundantVarDecl{
		file:      file,
		info:      info,
		onlyTypes: r.onlyTypes,
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}

	ast.Inspect(file.AST, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.FuncDecl:
			if n.Body != nil {
				w.checkFunc(n.Body)
			}
			return false
		case *ast.FuncLit:
			w.checkFunc(n.Body)
			return false
		}
		return true
	})

	return failures
}

// Name returns the rule name.
func (*RedundantVarDeclRule) Name() string {
	return "redundant-var-decl"
}

type lintRedundantVarDecl struct {
	file      *lint.File
	info      *types.Info
	onlyTypes map[string]struct{}
	onFailure func(lint.Failure)
}

// isSelectedType reports whether variables of type t are checked.
// A type is selected by its name qualified with the package name (e.g. "io.Reader")
// or with the full package path (e.g. "github.com/user/repo/pkg.Type").
func (w lintRedundantVarDecl) isSelectedType(t types.Type) bool {
	if len(w.onlyTypes) == 0 {
		return true
	}

	byName := types.TypeString(t, func(p *types.Package) string { return p.Name() })
	if _, ok := w.onlyTypes[byName]; ok {
		return true
	}

	_, ok := w.onlyTypes[types.TypeString(t, nil)]
	return ok
}

// checkFunc checks the statement lists of a function body.
// Nested function literals are checked on their own, because a goto cannot cross function boundaries.
func (w lintRedundantVarDecl) checkFunc(body *ast.BlockStmt) {
	// Moving a declaration down may make an existing goto jump over it, which does not compile.
	hasGoto := containsGoto(body)

	ast.Inspect(body, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.FuncLit:
			w.checkFunc(n.Body)
			return false
		case *ast.BlockStmt:
			if !hasGoto {
				w.checkStmts(n.List)
			}
		case *ast.CaseClause:
			if !hasGoto {
				w.checkStmts(n.Body)
			}
		case *ast.CommClause:
			if !hasGoto {
				w.checkStmts(n.Body)
			}
		}
		return true
	})
}

func (w lintRedundantVarDecl) checkStmts(stmts []ast.Stmt) {
	for i, stmt := range stmts {
		declStmt, ok := stmt.(*ast.DeclStmt)
		if !ok {
			continue
		}

		genDecl, ok := declStmt.Decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			continue
		}

		for j, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok || len(valueSpec.Values) > 0 {
				continue
			}

			for _, name := range valueSpec.Names {
				w.checkName(name, genDecl.Specs[j+1:], stmts[i+1:])
			}
		}
	}
}

// checkName reports name if the first statement referring to it is a short variable declaration
// that redeclares it with the very same type.
func (w lintRedundantVarDecl) checkName(name *ast.Ident, laterSpecs []ast.Spec, laterStmts []ast.Stmt) {
	if name.Name == "_" {
		return
	}

	obj, ok := w.info.Defs[name].(*types.Var)
	if !ok || !isValidType(obj.Type()) || !w.isSelectedType(obj.Type()) {
		return
	}

	for _, spec := range laterSpecs {
		if w.refersTo(spec, obj) {
			return
		}
	}

	for _, stmt := range laterStmts {
		assign, ok := stmt.(*ast.AssignStmt)
		if !ok || assign.Tok != token.DEFINE {
			if w.refersTo(stmt, obj) {
				return
			}
			continue
		}

		idx := w.lhsIndex(assign, obj)
		if idx < 0 {
			if w.refersTo(stmt, obj) {
				return
			}
			continue
		}

		for _, rhs := range assign.Rhs {
			if w.refersTo(rhs, obj) {
				return
			}
		}

		inferred := w.inferredType(assign, idx)
		if inferred == nil || !isValidType(inferred) || !types.Identical(obj.Type(), inferred) {
			return
		}

		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       name,
			Category:   lint.FailureCategoryStyle,
			Failure: fmt.Sprintf(
				"redundant declaration of %q: it is redeclared by := on line %d; remove the var declaration",
				name.Name, w.file.ToPosition(assign.Pos()).Line,
			),
		})
		return
	}
}

// lhsIndex returns the index of the identifier on the left side of assign that redeclares obj, or -1.
func (w lintRedundantVarDecl) lhsIndex(assign *ast.AssignStmt, obj types.Object) int {
	for i, lhs := range assign.Lhs {
		id, ok := lhs.(*ast.Ident)
		if ok && w.info.Uses[id] == obj {
			return i
		}
	}
	return -1
}

func (w lintRedundantVarDecl) refersTo(node ast.Node, obj types.Object) bool {
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		if found {
			return false
		}
		if id, ok := n.(*ast.Ident); ok && w.info.Uses[id] == obj {
			found = true
		}
		return !found
	})
	return found
}

// inferredType returns the type the variable at index idx of the left side of assign
// would get if the short variable declaration declared it.
func (w lintRedundantVarDecl) inferredType(assign *ast.AssignStmt, idx int) types.Type {
	if len(assign.Lhs) == len(assign.Rhs) {
		return w.standaloneType(assign.Rhs[idx])
	}

	if len(assign.Rhs) != 1 {
		return nil
	}

	rhs := ast.Unparen(assign.Rhs[0])
	if _, ok := rhs.(*ast.CallExpr); ok {
		// The results of a multi-value call are typed; they are recorded as a tuple.
		tuple, ok := w.info.TypeOf(rhs).(*types.Tuple)
		if !ok || idx >= tuple.Len() {
			return nil
		}
		return tuple.At(idx).Type()
	}

	// Comma-ok expressions: map index, type assertion and channel receive.
	switch idx {
	case 0:
		return w.standaloneType(rhs)
	case 1:
		return types.Typ[types.Bool]
	default:
		return nil
	}
}

// standaloneType returns the default type of expr evaluated outside of its assignment context.
// The type recorded by the type checker is not enough: in an assignment, untyped values
// (constants, comparisons, nil) are recorded with the type of the variable they are assigned to.
func (w lintRedundantVarDecl) standaloneType(expr ast.Expr) types.Type {
	if containsFuncLit(expr) {
		// Re-checking a function literal would add scopes to the shared package scope tree,
		// and the type of a function literal is never untyped anyway.
		return w.info.TypeOf(expr)
	}

	info := &types.Info{Types: map[ast.Expr]types.TypeAndValue{}}
	// The file set is only used to format error positions.
	err := types.CheckExpr(token.NewFileSet(), w.file.Pkg.TypesPkg(), expr.Pos(), expr, info)
	if err != nil {
		return nil
	}

	t := info.TypeOf(expr)
	if t == nil {
		return nil
	}
	return types.Default(t)
}

func containsFuncLit(expr ast.Expr) bool {
	found := false
	ast.Inspect(expr, func(n ast.Node) bool {
		if _, ok := n.(*ast.FuncLit); ok {
			found = true
		}
		return !found
	})
	return found
}

func isValidType(t types.Type) bool {
	return t != nil && t != types.Typ[types.Invalid]
}

func containsGoto(body *ast.BlockStmt) bool {
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		if found {
			return false
		}
		switch n := n.(type) {
		case *ast.FuncLit:
			return false
		case *ast.BranchStmt:
			if n.Tok == token.GOTO {
				found = true
			}
		}
		return !found
	})
	return found
}
