package fixtures

func compliantFunc(a, b int, c string) {} // Must not match - compliant with rule

func nonCompliantFunc1(a int, b int, c string) {} // MATCH /repeated argument type "int" can be omitted/
func nonCompliantFunc2(a int, b, c int)        {} // MATCH /repeated argument type "int" can be omitted/

type myStruct struct{}

func (m myStruct) compliantMethod(a, b int, c string) {} // Must not match - compliant with rule

func (m myStruct) nonCompliantMethod1(a int, b int, c string) {} // MATCH /repeated argument type "int" can be omitted/
func (m myStruct) nonCompliantMethod2(a int, b, c int)        {} // MATCH /repeated argument type "int" can be omitted/

func variadicFunction(a int, b ...int) {} // Must not match - variadic parameters are a special case

func singleArgFunction(a int) {} // Must not match - only one argument

func multiTypeArgs(a int, b string, c float64) {} // Must not match - different types for each argument

func mixedCompliance(a, b int, c int, d string) {} // MATCH /repeated argument type "int" can be omitted/ - 'c int' could be combined with 'a, b int'

func unnamedArgs(int, int, string) {} // Must not match - unnamed arguments cannot be grouped
