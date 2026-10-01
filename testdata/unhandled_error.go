package fixtures

import (
	"bytes"
	"fmt"
	fmt2 "fmt"
	"os"
	"strings"
)

func unhandledError1(a int) (int, error) {
	return a, nil
}

func unhandledError2() error {
	_, err := unhandledError1(1)
	unhandledError1(1)   // MATCH /Unhandled error in call to function unhandledError1/
	fmt.Fprintf(nil, "") // MATCH /Unhandled error in call to function fmt.Fprintf/

	var sb strings.Builder
	fmt.Fprintf(&sb, "formatted string: %v", 1)

	var bb bytes.Buffer
	fmt2.Fprintf(&bb, "formatted string: %v", 1)

	fmt.Fprintf(os.Stdout, "") // MATCH /Unhandled error in call to function fmt.Fprintf/
	os.Chdir("..")             // MATCH /Unhandled error in call to function os.Chdir/
	_ = os.Chdir("..")
	return err
}

type fixtureConcreteError struct{}

func (*fixtureConcreteError) Error() string { return "fixture" }

func ignoredConcreteError() *fixtureConcreteError { return nil }

func ignoredTupleConcreteError() (int, *fixtureConcreteError) { return 0, nil }

func unhandledConcreteError() {
	ignoredConcreteError()      // MATCH /Unhandled error in call to function ignoredConcreteError/
	ignoredTupleConcreteError() // MATCH /Unhandled error in call to function ignoredTupleConcreteError/
}
