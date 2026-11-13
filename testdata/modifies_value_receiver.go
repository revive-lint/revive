package fixtures

type data struct {
	num   int
	key   *string
	items map[string]bool
}

func (this data) vmethod() {
	this.num = 8 // MATCH /suspicious assignment to a by-value method receiver/
	*this.key = "v.key"
	this.items = make(map[string]bool) // MATCH /suspicious assignment to a by-value method receiver/
	this.items["vmethod"] = true
}

func (a A) Foo() *A {
	a.whatever = true
	return &a
}

func (a A) Clone() (*A, error) {
	a.whatever = true
	return &a, nil
}

// WithBin will set the specific bin path to the builder.
func (b JailerCommandBuilder) WithBin(bin string) JailerCommandBuilder {
	b.bin = bin
	return b
}

func (this data) incrementDecrement() {
	this.num++ // MATCH /suspicious assignment to a by-value method receiver/
	this.num-- // MATCH /suspicious assignment to a by-value method receiver/
	other++
}

func (this data) assignmentOperators() {
	this.num += 1  // MATCH /suspicious assignment to a by-value method receiver/
	this.num -= 1  // MATCH /suspicious assignment to a by-value method receiver/
	this.num *= 1  // MATCH /suspicious assignment to a by-value method receiver/
	this.num /= 1  // MATCH /suspicious assignment to a by-value method receiver/
	this.num %= 1  // MATCH /suspicious assignment to a by-value method receiver/
	this.num &= 1  // MATCH /suspicious assignment to a by-value method receiver/
	this.num ^= 1  // MATCH /suspicious assignment to a by-value method receiver/
	this.num |= 1  // MATCH /suspicious assignment to a by-value method receiver/
	this.num >>= 1 // MATCH /suspicious assignment to a by-value method receiver/
	this.num <<= 1 // MATCH /suspicious assignment to a by-value method receiver/
}

type counter int

func (c counter) increment() {
	c++ // MATCH /suspicious assignment to a by-value method receiver/
}

func (c counter) decrement() {
	c-- // MATCH /suspicious assignment to a by-value method receiver/
}
