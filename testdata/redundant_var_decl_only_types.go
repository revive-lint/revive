// Test of the onlyTypes option: only error, []byte, io.Reader and fixtures.flag variables are checked.

package fixtures

import (
	"io"
	"strings"
)

func readAll() ([]byte, error) { return nil, nil }

func newReader() (io.Reader, error) { return strings.NewReader(""), nil }

type flag bool

func parseFlag() (flag, error) { return true, nil }

func selectedError() error {
	var err error // MATCH /redundant declaration of "err": it is redeclared by := on line 20; remove the var declaration/
	data, err := readAll()
	_ = data
	return err
}

func selectedSlice() ([]byte, error) {
	var data []byte // MATCH /redundant declaration of "data": it is redeclared by := on line 27; remove the var declaration/
	data, err := readAll()
	return data, err
}

func selectedQualifiedType() (io.Reader, error) {
	var r io.Reader // MATCH /redundant declaration of "r": it is redeclared by := on line 33; remove the var declaration/
	r, err := newReader()
	return r, err
}

func selectedLocalType() (flag, error) {
	var f flag // MATCH /redundant declaration of "f": it is redeclared by := on line 39; remove the var declaration/
	f, err := parseFlag()
	return f, err
}

func notSelectedInt() int {
	var foo int
	foo, bar := 2, 3
	return foo + bar
}

func notSelectedBool(m map[string]int) (int, bool) {
	var ok bool
	v, ok := m["k"]
	return v, ok
}
