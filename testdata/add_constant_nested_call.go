package fixtures

func consume(any) {}

func addConstantNestedCall() {
	consume([]int{100}) // MATCH /avoid magic numbers like '100', create a named constant for it/
}
