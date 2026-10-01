package fixtures

import "fmt"

func unhandledErrorConfigResetOK() {
	fmt.Print("reported") // MATCH /Unhandled error in call to function fmt.Print/
}
