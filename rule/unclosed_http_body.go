package rule

import (
	"go/ast"
	"go/token"
	"go/types"
	"slices"

	"github.com/mgechev/revive/internal/astutils"
	"github.com/mgechev/revive/lint"
)

// UnclosedHTTPBodyRule spots HTTP responses whose body is never closed.
//
// The analysis is local to a function: a response returned by a call is considered handled if its body is closed
// or if the response (or its body) leaves the function, for example by being passed to another function, returned,
// or stored in a field. In the latter case the rule assumes that the receiver closes the body.
type UnclosedHTTPBodyRule struct{}

// Apply applies the rule to given file.
func (*UnclosedHTTPBodyRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if file.Pkg.TypeCheck() != nil {
		return nil
	}

	w := &lintUnclosedHTTPBody{info: file.Pkg.TypesInfo(), pkg: file.Pkg, nopBodyFuncs: map[*types.Func]bool{}}
	for _, decl := range file.AST.Decls {
		w.checkDecl(decl)
	}

	return w.failures
}

// Name returns the rule name.
func (*UnclosedHTTPBodyRule) Name() string {
	return "unclosed-http-body"
}

const unclosedHTTPBodyMessage = "response body must be closed"

type lintUnclosedHTTPBody struct {
	info     *types.Info
	pkg      *lint.Package
	failures []lint.Failure
	// decl is the top-level declaration being checked.
	decl ast.Decl
	// nopBodyFuncs caches, for functions of the package, if they replace the body of a response with a no-op closer.
	nopBodyFuncs map[*types.Func]bool
}

// httpResponseProducer is a call that returns an *http.Response bound to a local variable.
type httpResponseProducer struct {
	call *ast.CallExpr
	lhs  *ast.Ident
	// block is the innermost block (or case clause) containing the call.
	block   ast.Node
	handled bool
}

// httpResponseUse is a use of a variable holding an *http.Response.
type httpResponseUse struct {
	pos token.Pos
	// lit is the innermost function literal enclosing the use, if any.
	lit *ast.FuncLit
	// deferredLit is the innermost deferred function literal enclosing the use, if any.
	// Its body runs when the enclosing function returns.
	deferredLit *ast.FuncLit
}

func (w *lintUnclosedHTTPBody) checkDecl(decl ast.Decl) {
	w.decl = decl
	producers := map[types.Object][]*httpResponseProducer{}

	walkWithStack(decl, func(node ast.Node, stack []ast.Node) {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return
		}

		idx := w.responseIndex(call)
		if idx < 0 {
			return
		}

		obj, lhs, report := w.binding(call, idx, stack)
		switch {
		case report:
			w.addFailure(call)
		case obj != nil:
			producers[obj] = append(producers[obj], &httpResponseProducer{call: call, lhs: lhs, block: innermostBlock(stack)})
		}
	})

	for obj, prods := range producers {
		w.checkProducers(decl, obj, prods)
	}
}

// checkProducers reports the producers of the given variable whose response is not handled.
func (w *lintUnclosedHTTPBody) checkProducers(decl ast.Decl, obj types.Object, prods []*httpResponseProducer) {
	lhsIdents := map[*ast.Ident]bool{}
	for _, p := range prods {
		lhsIdents[p.lhs] = true
	}

	var uses []httpResponseUse
	var replacements []*ast.AssignStmt
	walkWithStack(decl, func(node ast.Node, stack []ast.Node) {
		if assign, ok := node.(*ast.AssignStmt); ok && w.replacesBodyWithNop(assign, obj) {
			replacements = append(replacements, assign)
			return
		}

		id, ok := node.(*ast.Ident)
		if !ok || lhsIdents[id] || w.info.Uses[id] != obj {
			return
		}

		if w.isResponseHandled(stack) {
			uses = append(uses, httpResponseUse{pos: id.Pos(), lit: innermostFuncLit(stack), deferredLit: innermostDeferredFuncLit(stack)})
		}
	})

	for _, use := range uses {
		for _, p := range prods {
			if useAppliesTo(use, p, prods) && !isBodyReplacedBefore(use, p, replacements) {
				p.handled = true
			}
		}
	}

	for _, p := range prods {
		if !p.handled {
			w.addFailure(p.call)
		}
	}
}

// useAppliesTo tells if the use of the variable may refer to the response returned by the producer.
//
// A use inside a function literal that does not contain the producer may run at any time (e.g. a deferred closure),
// so it applies to every producer. Otherwise it applies to a producer that precedes it,
// unless another producer between them always runs before the use, i.e. it is in a block that contains the use.
// Producers in other branches, like in
//
//	if cond {
//		resp, err = a()
//	} else {
//		resp, err = b()
//	}
//	defer resp.Body.Close()
//
// do not hide each other.
func useAppliesTo(use httpResponseUse, p *httpResponseProducer, prods []*httpResponseProducer) bool {
	if use.lit != nil && (p.call.Pos() < use.lit.Pos() || p.call.End() > use.lit.End()) {
		return true
	}

	if p.call.Pos() >= use.pos {
		return false
	}

	for _, other := range prods {
		between := other.call.Pos() > p.call.Pos() && other.call.Pos() < use.pos
		if between && other.block != nil && other.block.Pos() <= use.pos && use.pos < other.block.End() {
			return false
		}
	}

	return true
}

// isBodyReplacedBefore tells if the body of the response returned by the producer is replaced by a no-op closer
// before the use runs, like in
//
//	resp.Body = io.NopCloser(bytes.NewReader(b))
//	defer resp.Body.Close()
//
// In that case the use does not close the original body, nor makes it escape.
// A deferred function literal runs when the function returns, so after any replacement outside of it.
func isBodyReplacedBefore(use httpResponseUse, p *httpResponseProducer, replacements []*ast.AssignStmt) bool {
	for _, r := range replacements {
		if r.Pos() < p.call.End() {
			continue
		}

		if r.End() <= use.pos {
			return true
		}

		lit := use.deferredLit
		if lit != nil && !containsNode(lit, r) && !containsNode(lit, p.call) {
			return true
		}
	}

	return false
}

// replacesBodyWithNop tells if the assignment replaces the body of the response held by the variable
// with a body that does not need to be closed, like [io.NopCloser].
func (w *lintUnclosedHTTPBody) replacesBodyWithNop(assign *ast.AssignStmt, obj types.Object) bool {
	return slices.ContainsFunc(w.nopBodyTargets(assign), func(x ast.Expr) bool {
		id, ok := x.(*ast.Ident)
		return ok && w.info.Uses[id] == obj
	})
}

// nopBodyTargets returns the expressions X of the assignments of a body that does not need to be closed to X.Body
// in the given assignment statement.
func (w *lintUnclosedHTTPBody) nopBodyTargets(assign *ast.AssignStmt) []ast.Expr {
	if len(assign.Lhs) != len(assign.Rhs) {
		return nil
	}

	var targets []ast.Expr
	for i, lhs := range assign.Lhs {
		sel, ok := lhs.(*ast.SelectorExpr)
		if ok && sel.Sel.Name == "Body" && w.isNopBody(assign.Rhs[i]) {
			targets = append(targets, sel.X)
		}
	}

	return targets
}

// isNopBody tells if the expression is a body that does not need to be closed: a call to [io.NopCloser] or [http.NoBody].
func (w *lintUnclosedHTTPBody) isNopBody(expr ast.Expr) bool {
	switch e := ast.Unparen(expr).(type) {
	case *ast.CallExpr:
		fn := w.calleeFunc(e)
		if fn == nil {
			return false
		}
		name := fn.FullName()
		return name == "io.NopCloser" || name == "io/ioutil.NopCloser"
	case *ast.SelectorExpr:
		v, ok := w.info.Uses[e.Sel].(*types.Var)
		return ok && v.Pkg() != nil && v.Pkg().Path() == "net/http" && v.Name() == "NoBody"
	}

	return false
}

// calleeFunc returns the function or method called, or nil if it is not statically known.
func (w *lintUnclosedHTTPBody) calleeFunc(call *ast.CallExpr) *types.Func {
	var id *ast.Ident
	switch fun := ast.Unparen(call.Fun).(type) {
	case *ast.Ident:
		id = fun
	case *ast.SelectorExpr:
		id = fun.Sel
	case *ast.IndexExpr: // generic function instantiation
		id = identOf(fun.X)
	case *ast.IndexListExpr:
		id = identOf(fun.X)
	}

	if id == nil {
		return nil
	}

	fn, ok := w.info.Uses[id].(*types.Func)
	if !ok {
		return nil
	}

	return fn.Origin()
}

func identOf(expr ast.Expr) *ast.Ident {
	switch e := expr.(type) {
	case *ast.Ident:
		return e
	case *ast.SelectorExpr:
		return e.Sel
	}
	return nil
}

// returnsNopBody tells if the call is to a function of the package that replaces the body of a response
// with a body that does not need to be closed, like a helper that reads and closes the body, then returns
// the response with its body replaced by [io.NopCloser].
// Functions declared in other packages are not analyzed.
func (w *lintUnclosedHTTPBody) returnsNopBody(call *ast.CallExpr) bool {
	fn := w.calleeFunc(call)
	if fn == nil || fn.Pkg() != w.pkg.TypesPkg() {
		return false
	}

	if result, ok := w.nopBodyFuncs[fn]; ok {
		return result
	}

	result := false
	if decl := w.funcDecl(fn); decl != nil && decl.Body != nil {
		ast.Inspect(decl.Body, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.AssignStmt:
				result = slices.ContainsFunc(w.nopBodyTargets(n), func(x ast.Expr) bool {
					return isHTTPResponseOrValue(w.info.TypeOf(x))
				})
			case *ast.CompositeLit:
				result = w.isNopBodyResponseLit(n)
			}
			return !result
		})
	}

	w.nopBodyFuncs[fn] = result
	return result
}

// isNopBodyResponseLit tells if the composite literal is an [http.Response] with a body that does not need to be closed,
// like a call to [io.NopCloser].
func (w *lintUnclosedHTTPBody) isNopBodyResponseLit(lit *ast.CompositeLit) bool {
	if !isHTTPResponseOrValue(w.info.TypeOf(lit)) {
		return false
	}

	return slices.ContainsFunc(lit.Elts, func(elt ast.Expr) bool {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			return false
		}
		key, ok := kv.Key.(*ast.Ident)
		return ok && key.Name == "Body" && w.isNopBody(kv.Value)
	})
}

// funcDecl returns the declaration of the function in the package, or nil if it is not found.
func (w *lintUnclosedHTTPBody) funcDecl(fn *types.Func) *ast.FuncDecl {
	for _, file := range w.pkg.Files() {
		if file.AST == nil || file.AST.Pos() > fn.Pos() || fn.Pos() >= file.AST.End() {
			continue
		}
		for _, decl := range file.AST.Decls {
			if fd, ok := decl.(*ast.FuncDecl); ok && w.info.Defs[fd.Name] == fn {
				return fd
			}
		}
	}

	return nil
}

// responseIndex returns the position of the *http.Response in the results of the call, or -1 if there is none.
func (w *lintUnclosedHTTPBody) responseIndex(call *ast.CallExpr) int {
	fun, ok := w.info.Types[call.Fun]
	if !ok || fun.IsType() || fun.IsBuiltin() || w.isNoopBodyCall(call) || w.returnsNopBody(call) {
		return -1
	}

	switch t := w.info.TypeOf(call).(type) {
	case *types.Tuple:
		for i := range t.Len() {
			if isHTTPResponse(t.At(i).Type()) {
				return i
			}
		}
	default:
		if isHTTPResponse(t) {
			return 0
		}
	}

	return -1
}

// isNoopBodyCall tells if the call returns a response whose body does not need to be closed.
func (w *lintUnclosedHTTPBody) isNoopBodyCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	fn, ok := w.info.Uses[sel.Sel].(*types.Func)
	return ok && fn.FullName() == "(*net/http/httptest.ResponseRecorder).Result"
}

// binding returns the local variable the response of the call is assigned to, and its identifier.
// The boolean result is true when the response is discarded and must be reported.
// A nil object and a false boolean mean that the response escapes the analysis and is considered handled.
func (w *lintUnclosedHTTPBody) binding(call *ast.CallExpr, idx int, stack []ast.Node) (types.Object, *ast.Ident, bool) {
	var target ast.Expr
	switch parent := parentOf(stack).(type) {
	case *ast.ExprStmt, *ast.DeferStmt, *ast.GoStmt:
		return nil, nil, true
	case *ast.AssignStmt:
		target = assignTarget(parent.Lhs, parent.Rhs, call, idx)
	case *ast.ValueSpec:
		target = assignTarget(identsToExprs(parent.Names), parent.Values, call, idx)
	default:
		return nil, nil, false
	}

	id, ok := target.(*ast.Ident)
	if !ok {
		return nil, nil, false
	}

	if id.Name == "_" {
		return nil, nil, true
	}

	obj := w.info.ObjectOf(id)
	if obj == nil || obj.Pkg() == nil || obj.Parent() == obj.Pkg().Scope() {
		return nil, nil, false // package-level variables are not checked
	}

	return obj, id, false
}

// assignTarget returns the left-hand side expression that receives the response returned by the call.
func assignTarget(lhs, rhs []ast.Expr, call *ast.CallExpr, idx int) ast.Expr {
	if len(rhs) == 1 && len(lhs) > 1 {
		if idx < len(lhs) {
			return lhs[idx]
		}
		return nil
	}

	for i, r := range rhs {
		if r == call && i < len(lhs) {
			return lhs[i]
		}
	}

	return nil
}

// isResponseHandled tells if the use of a response variable (the last node of the stack) closes its body
// or makes the response escape the function.
func (w *lintUnclosedHTTPBody) isResponseHandled(stack []ast.Node) bool {
	id := stack[len(stack)-1]
	switch parent := parentOf(stack).(type) {
	case *ast.SelectorExpr:
		if w.isResponseWriteCall(stack[:len(stack)-1]) {
			return true // Response.Write closes the body after sending it
		}
		if parent.Sel.Name != "Body" {
			return false // other fields and methods, like resp.StatusCode
		}
		return w.isBodyHandled(stack[:len(stack)-1], map[types.Object]bool{})
	case *ast.CallExpr:
		if parent.Fun == id {
			return false
		}
		t := argParamType(w.info, parent, id)
		// passing the response to an interface parameter, like in fmt.Print(resp), does not close the body
		return t == nil || !types.IsInterface(t)
	case *ast.AssignStmt:
		return !isInExprs(id, parent.Lhs)
	case *ast.BinaryExpr, *ast.ParenExpr:
		return false
	}

	// returned, sent, stored in a composite literal, address taken, ...
	return true
}

// isResponseWriteCall tells if the selector (the last node of the stack) is a call to [http.Response.Write].
func (w *lintUnclosedHTTPBody) isResponseWriteCall(stack []ast.Node) bool {
	sel, ok := stack[len(stack)-1].(*ast.SelectorExpr)
	if !ok {
		return false
	}

	call, ok := parentOf(stack).(*ast.CallExpr)
	if !ok || call.Fun != sel {
		return false
	}

	fn, ok := w.info.Uses[sel.Sel].(*types.Func)
	return ok && fn.FullName() == "(*net/http.Response).Write"
}

// isBodyHandled tells if the use of a response body (the last node of the stack) closes it
// or makes it escape the function.
func (w *lintUnclosedHTTPBody) isBodyHandled(stack []ast.Node, visited map[types.Object]bool) bool {
	body := stack[len(stack)-1]
	switch parent := parentOf(stack).(type) {
	case *ast.SelectorExpr:
		return parent.Sel.Name == "Close"
	case *ast.CallExpr:
		if parent.Fun == body {
			return false
		}
		t := argParamType(w.info, parent, body)
		return t == nil || hasCloseMethod(t)
	case *ast.AssignStmt:
		if isInExprs(body, parent.Lhs) {
			return false
		}
		return w.isAliasHandled(parent.Lhs, parent.Rhs, body, visited)
	case *ast.ValueSpec:
		return w.isAliasHandled(identsToExprs(parent.Names), parent.Values, body, visited)
	case *ast.BinaryExpr, *ast.ParenExpr:
		return false
	}

	// returned, sent, stored in a composite literal, ...
	return true
}

// isAliasHandled checks the variable that receives the body in an assignment or declaration, like in `body := resp.Body`.
func (w *lintUnclosedHTTPBody) isAliasHandled(lhs, rhs []ast.Expr, body ast.Node, visited map[types.Object]bool) bool {
	var target ast.Expr
	for i, r := range rhs {
		if r == body && i < len(lhs) {
			target = lhs[i]
		}
	}

	id, ok := target.(*ast.Ident)
	if !ok {
		return true // stored in a field, a map, ...
	}

	if id.Name == "_" {
		return false
	}

	obj := w.info.ObjectOf(id)
	if obj == nil || obj.Pkg() == nil || obj.Parent() == obj.Pkg().Scope() {
		return true // stored in a package-level variable
	}

	if visited[obj] {
		return false
	}
	visited[obj] = true

	return w.isAliasUsedToClose(obj, visited)
}

// isAliasUsedToClose tells if a variable holding a response body is closed or escapes the function.
func (w *lintUnclosedHTTPBody) isAliasUsedToClose(alias types.Object, visited map[types.Object]bool) bool {
	handled := false
	walkWithStack(w.decl, func(node ast.Node, stack []ast.Node) {
		if handled {
			return
		}
		id, ok := node.(*ast.Ident)
		if ok && w.info.Uses[id] == alias && w.isBodyHandled(stack, visited) {
			handled = true
		}
	})

	return handled
}

// argParamType returns the type of the parameter that receives the argument, or nil if it cannot be determined.
func argParamType(info *types.Info, call *ast.CallExpr, arg ast.Node) types.Type {
	sig, ok := types.Unalias(info.TypeOf(call.Fun)).(*types.Signature)
	if !ok {
		return nil
	}

	idx := -1
	for i, a := range call.Args {
		if a == arg {
			idx = i
			break
		}
	}

	params := sig.Params()
	if idx < 0 || params.Len() == 0 {
		return nil
	}

	if sig.Variadic() && idx >= params.Len()-1 {
		last := params.At(params.Len() - 1).Type()
		if call.Ellipsis.IsValid() {
			return last
		}
		if s, ok := last.(*types.Slice); ok {
			return s.Elem()
		}
		return nil
	}

	if idx >= params.Len() {
		return nil
	}

	return params.At(idx).Type()
}

func hasCloseMethod(t types.Type) bool {
	obj, _, _ := types.LookupFieldOrMethod(t, true, nil, "Close")
	_, ok := obj.(*types.Func)
	return ok
}

func isHTTPResponse(t types.Type) bool {
	return astutils.IsPointerToPkgDotType(types.Unalias(t), "net/http", "Response")
}

// isHTTPResponseOrValue tells if the type is http.Response or *http.Response.
func isHTTPResponseOrValue(t types.Type) bool {
	return t != nil && (isHTTPResponse(t) || isHTTPResponse(types.NewPointer(types.Unalias(t))))
}

func identsToExprs(idents []*ast.Ident) []ast.Expr {
	exprs := make([]ast.Expr, len(idents))
	for i, id := range idents {
		exprs[i] = id
	}
	return exprs
}

func isInExprs(node ast.Node, exprs []ast.Expr) bool {
	return slices.ContainsFunc(exprs, func(e ast.Expr) bool { return e == node })
}

func parentOf(stack []ast.Node) ast.Node {
	if len(stack) < 2 {
		return nil
	}
	return stack[len(stack)-2]
}

// innermostBlock returns the innermost block statement or case clause of the stack.
func innermostBlock(stack []ast.Node) ast.Node {
	for _, node := range slices.Backward(stack) {
		switch node.(type) {
		case *ast.BlockStmt, *ast.CaseClause, *ast.CommClause:
			return node
		}
	}
	return nil
}

// innermostDeferredFuncLit returns the innermost deferred function literal of the stack, like in `defer func() { ... }()`.
func innermostDeferredFuncLit(stack []ast.Node) *ast.FuncLit {
	for i := len(stack) - 1; i >= 2; i-- {
		lit, ok := stack[i].(*ast.FuncLit)
		if !ok {
			continue
		}
		call, ok := stack[i-1].(*ast.CallExpr)
		if !ok || call.Fun != lit {
			continue
		}
		if _, ok := stack[i-2].(*ast.DeferStmt); ok {
			return lit
		}
	}
	return nil
}

func containsNode(outer, inner ast.Node) bool {
	return outer.Pos() <= inner.Pos() && inner.End() <= outer.End()
}

func innermostFuncLit(stack []ast.Node) *ast.FuncLit {
	for _, node := range slices.Backward(stack) {
		if lit, ok := node.(*ast.FuncLit); ok {
			return lit
		}
	}
	return nil
}

// walkWithStack calls fn for every node under root, along with the path from root to the node (the node included).
func walkWithStack(root ast.Node, fn func(node ast.Node, stack []ast.Node)) {
	var stack []ast.Node
	ast.Inspect(root, func(node ast.Node) bool {
		if node == nil {
			stack = stack[:len(stack)-1]
			return true
		}
		stack = append(stack, node)
		fn(node, stack)
		return true
	})
}

func (w *lintUnclosedHTTPBody) addFailure(call *ast.CallExpr) {
	w.failures = append(w.failures, lint.Failure{
		Category:   lint.FailureCategoryBadPractice,
		Confidence: 0.8,
		Node:       call,
		Failure:    unclosedHTTPBodyMessage,
	})
}
