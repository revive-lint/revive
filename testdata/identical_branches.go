package fixtures

func identicalBranches() {
	if true { // MATCH /both branches of the if are identical/

	} else {

	}

	if true {

	}

	if true {
		print()
	} else {
	}

	if true { // MATCH /both branches of the if are identical/
		print()
	} else {
		print()
	}

	if true {
		if true { // MATCH /both branches of the if are identical/
			print()
		} else {
			print()
		}
	} else {
		println()
	}

	if true {
		println("something")
	} else {
		println("else")
	}
}

var packageLevelFuncLit = func() {
	if true { // MATCH /both branches of the if are identical/
		print()
	} else {
		print()
	}
}

var packageLevelFuncLitInComposite = []func(){
	func() {
		if true { // MATCH /both branches of the if are identical/
			print()
		} else {
			print()
		}
	},
}
