package fixtures

func identicalSwitchBranches() {
	switch a { // MATCH /"switch" with identical branches (lines 6 and 10)/
	// expected values
	case 1:
		foo()
	case 2:
		bar()
	case 3:
		foo()
	default:
		return newError("blah")
	}

	// MATCH:18 /"switch" with identical branches (lines 20 and 24)/
	// MATCH:18 /"switch" with identical branches (lines 22 and 26)/
	switch a {
	// expected values
	case 1:
		foo()
	case 2:
		bar()
	case 3:
		foo()
	default:
		bar()
	}

	switch a { // MATCH /"switch" with identical branches (lines 32 and 34)/
	// expected values
	case 1:
		foo()
	case 3:
		foo()
	default:

	}

	// Skip untagged switch
	switch {
	case a > b:
		foo()
	default:
		foo()
	}

	// Do not warn on fallthrough

	switch a {
	case 1:
		foo()
		fallthrough
	case 2:
		fallthrough
	case 3:
		foo()
	case 4:
		fallthrough
	default:
		bar()
	}

	// skip type switch
	switch v := value.(type) {
	case int:
		println("dup", v)
	case string:
		println("dup", v)
	case bool:
		println("dup", v)
	case float64:
		println("dup", v)
	default:
		println("dup", v)
	}
}

var packageLevelFuncLit = func() {
	switch a { // MATCH /"switch" with identical branches (lines 81 and 83)/
	case 1:
		foo()
	case 2:
		foo()
	}
}
