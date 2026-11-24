package fixtures

// Type checking must not panic on nested indexed composite literals with invalid keys and values (issue #59).

var a = [][]int{
	2: {1: 1, "x": 2, 3: "y"},
	{-1: 0, 1: undefined},
}

var b = map[string][]int{
	"a": {0: 1, 0: 2, 1 << 70: 3},
	1:   {2: 3},
}

var c = [...][2][]struct{ F int }{
	1: {1: {{F: "x"}, {G: 1}}},
	5: {5: {{}}},
}
