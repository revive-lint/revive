package fixtures

import (
	"errors"
	"os"
)

type tt int

func (t tt) m() {}

func deferrer1() {
	for {
		go func() {
			defer println()
		}()
		defer func() {}()
	}

	defer tt.m() // MATCH /be careful when deferring calls to methods without pointer receiver/

	defer func() error {
		return errors.New("error") //MATCH /return in a defer function has no effect/
	}()

	defer recover()

	recover() //MATCH /recover must be called inside a deferred function/

	defer deferrer()
}

func exitAfterDefer() {
	defer println("cleanup")
	os.Exit(1) // MATCH /os.Exit after a defer statement prevents deferred calls from running/
}
