// Package pkg ...
package pkg

type Foo struct {
}

type Bar struct {
}

type Baz struct {
}

type Reader interface {
	Read() error
}

type NamedInt int

type Callback func()
