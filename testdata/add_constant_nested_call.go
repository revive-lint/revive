package fixtures

const zero = 0

func consume(any) {}

func factory(int) func() int { return nil }

type doer struct{}

func (doer) Do() int { return zero }

func get([]int) doer { return doer{} }

func addConstantNestedCall() {
	consume([]int{100}) // MATCH /avoid magic numbers like '100', create a named constant for it/

	consume(func(_ [101]int) {}) // MATCH /avoid magic numbers like '101', create a named constant for it/

	consume([]any{factory(102)()}) // MATCH /avoid magic numbers like '102', create a named constant for it/

	consume(get([]int{103}).Do()) // MATCH /avoid magic numbers like '103', create a named constant for it/

	consume(map[string]int{"key": 104}) // MATCH /avoid magic numbers like '104', create a named constant for it/

	consume(struct {
		F int `json:"f"`
	}{F: 105}) // MATCH /avoid magic numbers like '105', create a named constant for it/
}
