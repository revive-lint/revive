package fixtures

func (this data) vmethod() {
	nil := true // MATCH /assignment creates a shadow of built-in identifier nil/
	iota = 1    // MATCH /assignment modifies built-in identifier iota/
}

func append(i, j int) { // MATCH /redefinition of the built-in function append/

}

type string int16 // MATCH /redefinition of the built-in type string/

func delete(set []int64, i int) (y []int64) { // MATCH /redefinition of the built-in function delete/
	for j, v := range set {
		if j != i {
			y = append(y, v)
		}
	}
	return
}

type any int // MATCH /redefinition of the built-in type any/

func any() {} // MATCH /redefinition of the built-in type any/

var any int // MATCH /redefinition of the built-in type any/

const any = 1 // MATCH /redefinition of the built-in type any/

var i, copy int // MATCH /redefinition of the built-in function copy/

// issue #792
type ()

func foo() {
	clear := 0 // Shall not match /redefinition of the built-in function clear/
	max := 0   // Shall not match /redefinition of the built-in function max/
	min := 0   // Shall not match /redefinition of the built-in function min/
	_ = clear
	_ = max
	_ = min
}

func foo1(new int) { // MATCH /redefinition of the built-in function new/
	_ = new
}

func foo2() (new int) { // MATCH /redefinition of the built-in function new/
	return
}

func foo3[new any]() { // MATCH /redefinition of the built-in function new/
}

type comparable int // MATCH /redefinition of the built-in type comparable/

func comparable() {} // MATCH /redefinition of the built-in type comparable/

var comparable int // MATCH /redefinition of the built-in type comparable/

const comparable = 1 // MATCH /redefinition of the built-in type comparable/

func foo4[comparable any]() { // MATCH /redefinition of the built-in type comparable/
}

type ( // MATCH /redefinition of the built-in function len/
	Local int
	len   int
)
