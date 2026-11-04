package fixtures

import (
	"bytes"
	"errors"
	"io"
	"os"
)

func call() ([]int, error) { return nil, nil }

func callMyErr() ([]int, *myErr) { return nil, nil }

type myErr struct{}

func (*myErr) Error() string { return "" }

func process(error) {}

// Reported cases.

func issueExample() error {
	var err error // MATCH /redundant declaration of "err": it is redeclared by := on line 24; remove the var declaration/
	insts, err := call()
	_ = insts
	return err
}

func unrelatedStatementsInBetween() error {
	var err error // MATCH /redundant declaration of "err": it is redeclared by := on line 34; remove the var declaration/

	name := "a"
	_ = name
	insts, err := call()
	_ = insts
	return err
}

func nonErrorType() int {
	var foo int // MATCH /redundant declaration of "foo": it is redeclared by := on line 41; remove the var declaration/
	foo, bar := 2, 3
	return foo + bar
}

func groupedDecl() (int, error) {
	var (
		err error // MATCH /redundant declaration of "err": it is redeclared by := on line 50; remove the var declaration/
		n   int
	)
	insts, err := call()
	_ = insts
	n = 1
	return n, err
}

func multipleNames() (int, int) {
	var a, b int // MATCH /redundant declaration of "a": it is redeclared by := on line 59; remove the var declaration/
	b = 1
	a, c := 2, 3
	return a + b, c
}

func commaOk(m map[string]int) (int, bool) {
	var ok bool // MATCH /redundant declaration of "ok": it is redeclared by := on line 65; remove the var declaration/
	v, ok := m["k"]
	return v, ok
}

func insideCaseClause(x int) error {
	switch x {
	case 1:
		var err error // MATCH /redundant declaration of "err": it is redeclared by := on line 73; remove the var declaration/
		insts, err := call()
		_ = insts
		return err
	}
	return nil
}

func insideFuncLit() {
	f := func() error {
		var err error // MATCH /redundant declaration of "err": it is redeclared by := on line 83; remove the var declaration/
		insts, err := call()
		_ = insts
		return err
	}
	_ = f
}

// Not reported cases.

func usedInBetween() error {
	var err error
	process(err)
	insts, err := call()
	_ = insts
	return err
}

func assignedInBetween() error {
	var err error
	err = errors.New("x")
	process(err)
	insts, err := call()
	_ = insts
	return err
}

func addressTakenInBetween() error {
	var err error
	p := &err
	_ = p
	insts, err := call()
	_ = insts
	return err
}

func capturedByDefer() (ret error) {
	var err error
	defer func() { ret = err }()
	insts, err := call()
	_ = insts
	return err
}

func usedInRHS() error {
	var err error
	insts, err := wrap(err)
	_ = insts
	return err
}

func wrap(err error) (int, error) { return 0, err }

func usedInLaterSpec() error {
	var (
		err  error
		copy = err
	)
	insts, err := call()
	_, _ = insts, copy
	return err
}

func plainAssignment() error {
	var err error
	_, err = call()
	return err
}

func shadowedInIfInit() error {
	var err error
	if insts, err := call(); err != nil {
		_ = insts
		return err
	}
	return err
}

func shadowedInNestedBlock() error {
	var err error
	{
		insts, err := call()
		_, _ = insts, err
	}
	return err
}

func shadowedInLoop() error {
	var err error
	for range 3 {
		insts, err := call()
		_, _ = insts, err
	}
	return err
}

func withInitializer() error {
	var err = errors.New("x")
	insts, err := call()
	_ = insts
	return err
}

func concreteErrorType() error {
	var err error
	insts, err := callMyErr()
	_ = insts
	return err
}

func interfaceToConcrete() error {
	var r io.Reader
	r, err := os.Open("x")
	if err != nil {
		return err
	}
	r = bytes.NewReader(nil)
	_ = r
	return nil
}

func untypedConstant() int64 {
	var n int64
	n, ok := 2, true
	_ = ok
	return n
}

func withGoto() error {
	var err error
	goto end
end:
	insts, err := call()
	_ = insts
	return err
}

func untypedNil() *int {
	var p *int
	p, q := nil, 1
	_ = q
	return p
}

type flag bool

func untypedBoolToNamedType(a, b int) flag {
	var f flag
	f, n := a == b, 1
	_ = n
	return f
}

func typedConstant() int64 {
	const c int64 = 2
	var n int64 // MATCH /redundant declaration of "n": it is redeclared by := on line 238; remove the var declaration/
	n, ok := c, true
	_ = ok
	return n
}

func namedBoolCommaOk(m map[string]int) flag {
	var ok flag
	v, ok := m["k"]
	_ = v
	return ok
}
