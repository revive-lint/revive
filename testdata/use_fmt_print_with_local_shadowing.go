package fixtures

import (
	"log"
)

func useFmtPrintShadowedByParam(println func(string)) {
	println("just testing")
}

func useFmtPrintShadowedByLocalVar() {
	print := log.Print
	print("just testing")
}

func useFmtPrintShadowedByNamedResult() (println func(string)) {
	println = func(string) {}
	println("just testing")
	return println
}

func useFmtPrintBuiltin() {
	println("just testing") // MATCH /avoid using built-in function "println", replace it by "fmt.Fprintln(os.Stderr, "just testing")"/
	print("just testing")   // MATCH /avoid using built-in function "print", replace it by "fmt.Fprint(os.Stderr, "just testing")"/
}
