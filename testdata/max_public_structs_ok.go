// Package pkg ...
package pkg

type Foo struct {
}

type Bar struct {
}

type Baz struct {
}

// Unexported names, even non-ASCII or underscore-prefixed ones, are not public.
type éclair struct {
}

type структура struct {
}

type αlpha struct {
}

type _foo struct {
}
