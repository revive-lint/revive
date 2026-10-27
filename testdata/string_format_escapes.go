package fixtures

import "fmt"

func escapes() {
	fmt.Errorf("trailing newline\n") // MATCH /must not end with a newline/
	fmt.Errorf("no trailing newline")
	fmt.Errorf(`raw string keeps \n as written`)
}
