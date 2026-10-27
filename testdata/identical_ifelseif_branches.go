package fixtures

func identicalIfElseIfBranches() {

	if true { // MATCH /"if...else if" chain with identical branches (lines 5 and 13)/
		print("something")
	} else if true {
		print("something else")
	} else if true {
		print("other thing")
	} else if false {
		println()
	} else {
		print("something")
	}

	if true { // MATCH /"if...else if" chain with identical branches (lines 17 and 23)/
		print("something")
	} else if true {
		print("something else")
	} else if true {
		print("other thing")
	} else if false {
		print("something")
	} else {
		println()
	}

	if true {
		print("something")
	} else if true {
		print("something else")
		if true { // MATCH /"if...else if" chain with identical branches (lines 33 and 35)/
			print("something")
		} else if false {
			print("something")
		} else {
			println()
		}
	}

	// Should not warn because even if branches are identical, the err variable is not
	if err := something(); err != nil {
		println(err)
	} else if err := somethingElse(); err != nil {
		println(err)
	}

	// Multiple identical pair of branches
	if a {
		foo()
	} else if b {
		bar()
	} else if c {
		foo()
	} else if d {
		bar()
	}
	// MATCH:50 /"if...else if" chain with identical branches (lines 50 and 54)/
	// MATCH:50 /"if...else if" chain with identical branches (lines 52 and 56)/

	if createFile() { // json:{"MATCH": "\"if...else if\" chain with identical branches (lines 62 and 66)","Confidence": 0.8}
		doSomething()
	} else if !delete() {
		return new("cannot delete file")
	} else if createFile() {
		doSomething()
	} else {
		return new("file error")
	}

	// Test confidence is reset
	if a { // json:{"MATCH": "\"if...else if\" chain with identical branches (lines 73 and 75)","Confidence": 1}
		foo()
	} else if b {
		foo()
	} else {
		bar()
	}
}

var packageLevelFuncLit = func() {
	if a { // MATCH /"if...else if" chain with identical branches (lines 83 and 85)/
		foo()
	} else if b {
		foo()
	}
}
