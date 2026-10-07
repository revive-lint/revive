// Configured with preserveScope, so moving the declaration of x is not suggested.

package fixtures

func superfluousElseConfigReset() {
	for {
		if x := fn1(); x != nil {
			continue
		} else {
			fn2()
		}
		x := fn2()
		fn3(x)
	}
}
