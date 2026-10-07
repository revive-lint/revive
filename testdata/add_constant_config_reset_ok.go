package fixtures

import "fmt"

func addConstantConfigResetOK() {
	fmt.Print(1) // MATCH /avoid magic numbers like '1', create a named constant for it/
}
