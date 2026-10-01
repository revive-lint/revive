package rule

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/mgechev/revive/internal/astutils"
	"github.com/mgechev/revive/lint"
)

// ConstantLogicalExprRule warns on constant logical expressions.
type ConstantLogicalExprRule struct{}

// Apply applies the rule to given file.
func (*ConstantLogicalExprRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	if err := file.Pkg.TypeCheck(); err != nil && file.Pkg.TypesInfo() == nil {
		return nil
	}
	w := &lintConstantLogicalExpr{file: file, onFailure: onFailure}
	ast.Walk(w, file.AST)
	return failures
}

// Name returns the rule name.
func (*ConstantLogicalExprRule) Name() string {
	return "constant-logical-expr"
}

type lintConstantLogicalExpr struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

func (w *lintConstantLogicalExpr) Visit(node ast.Node) ast.Visitor {
	if n, ok := node.(*ast.BinaryExpr); ok {
		if !w.isOperatorWithLogicalResult(n.Op) {
			return w
		}
		info := w.file.Pkg.TypesInfo()
		if hasUnstableEvaluation(n.X, info) || hasUnstableEvaluation(n.Y, info) {
			return w
		}

		subExpressionsAreNotEqual := astutils.GoFmt(n.X) != astutils.GoFmt(n.Y)
		if subExpressionsAreNotEqual {
			return w // nothing to say
		}
		if w.isNaNComparison(n) {
			return w
		}

		// Handles cases like: a <= a, a == a, a >= a
		if w.isEqualityOperator(n.Op) {
			w.newFailure(n, "expression always evaluates to true")
			return w
		}

		// Handles cases like: a < a, a > a, a != a
		if w.isInequalityOperator(n.Op) {
			w.newFailure(n, "expression always evaluates to false")
			return w
		}

		w.newFailure(n, "left and right hand-side sub-expressions are the same")
	}

	return w
}

func hasUnstableEvaluation(expr ast.Expr, info *types.Info) bool {
	unstable := false
	ast.Inspect(expr, func(node ast.Node) bool {
		if unstable || node == nil {
			return !unstable
		}
		switch n := node.(type) {
		case *ast.CallExpr:
			isTypeConversion := info != nil && info.Types[n.Fun].IsType()
			unstable = !isTypeConversion
		case *ast.UnaryExpr:
			unstable = n.Op == token.ARROW
		}
		return !unstable
	})
	return unstable
}

func (w *lintConstantLogicalExpr) isNaNComparison(expr *ast.BinaryExpr) bool {
	switch expr.Op {
	case token.EQL, token.NEQ, token.LEQ, token.GEQ:
	default:
		return false
	}

	typ := w.file.Pkg.TypeOf(expr.X)
	if typ == nil {
		return true // do not make a constant claim without type information
	}

	info := w.file.Pkg.TypesInfo()
	if info != nil && info.Types[expr.X].Value != nil {
		return false // constants cannot be NaN
	}

	return typeMayContainNaN(typ)
}

func typeMayContainNaN(typ types.Type) bool {
	switch t := typ.(type) {
	case *types.TypeParam, *types.Interface:
		return true
	case *types.Array:
		return typeMayContainNaN(t.Elem())
	case *types.Struct:
		for field := range t.Fields() {
			if typeMayContainNaN(field.Type()) {
				return true
			}
		}
		return false
	}

	underlying := typ.Underlying()
	if underlying != typ {
		return typeMayContainNaN(underlying)
	}
	if basic, ok := underlying.(*types.Basic); ok {
		switch basic.Kind() {
		case types.Float32, types.Float64, types.Complex64, types.Complex128:
			return true
		}
	}
	return false
}

func (*lintConstantLogicalExpr) isOperatorWithLogicalResult(t token.Token) bool {
	switch t {
	case token.LAND, token.LOR, token.EQL, token.LSS, token.GTR, token.NEQ, token.LEQ, token.GEQ:
		return true
	}

	return false
}

func (*lintConstantLogicalExpr) isEqualityOperator(t token.Token) bool {
	switch t {
	case token.EQL, token.LEQ, token.GEQ:
		return true
	}

	return false
}

func (*lintConstantLogicalExpr) isInequalityOperator(t token.Token) bool {
	switch t {
	case token.LSS, token.GTR, token.NEQ:
		return true
	}

	return false
}

func (w *lintConstantLogicalExpr) newFailure(node ast.Node, msg string) {
	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       node,
		Category:   lint.FailureCategoryLogic,
		Failure:    msg,
	})
}
