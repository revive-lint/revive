package fixtures

// new accepts a value to point to only since Go 1.26: nothing is reported here.

func useNewBeforeGo126(b []byte) *string {
	s := string(b)

	return &s
}
