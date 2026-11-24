// Configured with preserveScope, so moving the declaration of y is not suggested.

package fixtures

func indentErrorFlowConfigReset() {
	if cond {
		var x = fn2()
		fn3(x)
		return
	} else {
		y := fn2()
		fn3(y)
	}
	y := fn2()
	fn3(y)
}
