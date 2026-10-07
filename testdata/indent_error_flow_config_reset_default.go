// Same code as indent_error_flow_config_reset.go, but linted with the default configuration.

package fixtures

func indentErrorFlowConfigResetDefault() {
	if cond {
		var x = fn2()
		fn3(x)
		return
	} else { // MATCH /if block ends with a return statement, so drop this else and outdent its block/
		y := fn2()
		fn3(y)
	}
	y := fn2()
	fn3(y)
}
