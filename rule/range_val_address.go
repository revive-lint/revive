package rule

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/mgechev/revive/lint"
)

// RangeValAddress warns if address of range value is used dangerously.
type RangeValAddress struct{}

// Apply applies the rule to given file.
func (*RangeValAddress) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	if file.Pkg.IsAtLeastGoVersion(lint.Go122) {
		return failures
	}

	walker := rangeValAddress{
		file: file,
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}

	file.Pkg.TypeCheck() //nolint:errcheck,gosec // G104: Errors unhandled
	ast.Walk(walker, file.AST)

	return failures
}

// Name returns the rule name.
func (*RangeValAddress) Name() string {
	return "range-val-address"
}

type rangeValAddress struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

func (w rangeValAddress) Visit(node ast.Node) ast.Visitor {
	n, ok := node.(*ast.RangeStmt)
	if !ok {
		return w
	}

	value, ok := n.Value.(*ast.Ident)
	if !ok {
		return w
	}

	info := w.file.Pkg.TypesInfo()
	valueObj := info.ObjectOf(value)
	if valueObj == nil {
		return w
	}

	ast.Walk(rangeBodyVisitor{
		info:      info,
		valueObj:  valueObj,
		onFailure: w.onFailure,
	}, n.Body)

	return w
}

type rangeBodyVisitor struct {
	info      *types.Info
	valueObj  types.Object
	onFailure func(lint.Failure)
}

func (bw rangeBodyVisitor) Visit(node ast.Node) ast.Visitor {
	asgmt, ok := node.(*ast.AssignStmt)
	if !ok {
		return bw
	}

	for _, exp := range asgmt.Lhs {
		e, ok := exp.(*ast.IndexExpr)
		if !ok {
			continue
		}
		if bw.isAccessingRangeValueAddress(e.Index) { // e.g. a[&value]...
			bw.onFailure(bw.newFailure(e.Index))
		}
	}

	for _, exp := range asgmt.Rhs {
		switch e := exp.(type) {
		case *ast.UnaryExpr: // e.g. ...&value, ...&value.id
			if bw.isAccessingRangeValueAddress(e) {
				bw.onFailure(bw.newFailure(e))
			}
		case *ast.CallExpr:
			if fun, ok := e.Fun.(*ast.Ident); ok && fun.Name == "append" { // e.g. ...append(arr, &value)
				for _, v := range e.Args {
					if lit, ok := v.(*ast.CompositeLit); ok { // e.g. ...append(arr, v{id:&value})
						bw.checkCompositeLit(lit)
						continue
					}
					if bw.isAccessingRangeValueAddress(v) { // e.g. ...append(arr, &value)
						bw.onFailure(bw.newFailure(v))
					}
				}
			}
		case *ast.CompositeLit: // e.g. ...v{id:&value}
			bw.checkCompositeLit(e)
		}
	}
	return bw
}

func (bw rangeBodyVisitor) checkCompositeLit(comp *ast.CompositeLit) {
	for _, exp := range comp.Elts {
		e, ok := exp.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		if bw.isAccessingRangeValueAddress(e.Value) {
			bw.onFailure(bw.newFailure(e.Value))
		}
	}
}

func (bw rangeBodyVisitor) isAccessingRangeValueAddress(exp ast.Expr) bool {
	u, ok := exp.(*ast.UnaryExpr)
	if !ok || u.Op != token.AND {
		return false
	}

	return bw.isStoredInRangeValue(u.X)
}

// isStoredInRangeValue reports whether expr denotes the range value itself or memory stored inline in it,
// such as a struct field or an array element, as opposed to memory reached through a pointer.
func (bw rangeBodyVisitor) isStoredInRangeValue(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.Ident:
		return bw.info.Uses[e] == bw.valueObj
	case *ast.ParenExpr:
		return bw.isStoredInRangeValue(e.X)
	case *ast.SelectorExpr:
		sel, ok := bw.info.Selections[e]
		if !ok || sel.Kind() != types.FieldVal || sel.Indirect() {
			return false
		}
		return bw.isStoredInRangeValue(e.X)
	case *ast.IndexExpr:
		t := bw.info.TypeOf(e.X)
		if t == nil {
			return false
		}
		if _, ok := t.Underlying().(*types.Array); !ok {
			return false
		}
		return bw.isStoredInRangeValue(e.X)
	default:
		return false
	}
}

func (bw rangeBodyVisitor) newFailure(node ast.Node) lint.Failure {
	return lint.Failure{
		Category:   lint.FailureCategoryLogic,
		Node:       node,
		Confidence: 1,
		Failure:    fmt.Sprintf("suspicious assignment of '%s'. range-loop variables always have the same address", bw.valueObj.Name()),
	}
}
