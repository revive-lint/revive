// Same code as superfluous_else_config_reset.go, but linted with the default configuration.

package fixtures

func superfluousElseConfigResetDefault() {
	for {
		if x := fn1(); x != nil {
			continue
		} else { // MATCH /if block ends with a continue statement, so drop this else and outdent its block (move short variable declaration to its own line if necessary)/
			fn2()
		}
		x := fn2()
		fn3(x)
	}
}
