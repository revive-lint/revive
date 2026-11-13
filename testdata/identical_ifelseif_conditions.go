package fixtures

func identicalBranches() {
	// no failure for nested ifs
	if true {
		if true {
		}
	} else {
		if true {
		}
	}

	// single failure
	if a > 0 {
		print("something")
	} else if a < 0 {
		print("something else")
	} else if a == 0 {
		print("other thing")
	} else if a > 0 { // MATCH /"if...else if" chain with identical conditions (lines 14 and 20)/
		println()
	} else {
		print("something")
	}

	// multiple failures in the same if...else if chain
	if a > 0 {
		print("something")
	} else if a < 0 {
		print("something else")
	} else if a == 0 {
		print("other thing")
	} else if a > 0 { // MATCH /"if...else if" chain with identical conditions (lines 27 and 33)/
		println()
	} else if a == 0 { // MATCH /"if...else if" chain with identical conditions (lines 31 and 35)/
		print("other thing")
	} else {
		print("something")
	}

	// failures in nested if...else if
	if true {
		if false {
		} else if false { // MATCH /"if...else if" chain with identical conditions (lines 43 and 44)/
		}
	} else if foo() {

	} else {
		if false {
		} else if false { // MATCH /"if...else if" chain with identical conditions (lines 49 and 50)/
		}
	}
}

var packageLevelFuncLit = func() {
	if a {
	} else if a { // MATCH /"if...else if" chain with identical conditions (lines 56 and 57)/
	}
}
