package fixtures

import "fmt"

// Comparing a float with itself is also a standard way to detect NaN.
func isNaN(f float64) bool { return f != f }

type floatBox struct{ value float64 }

func sameFloatBox(value floatBox) bool { return value == value }

func sameInterface(value any) bool { return value == value }

func skip(f float64) bool { return f != g }

func foo1(f float64) bool { return foo2(2.) > foo2(2.) }

func foo2(f float64) bool { return f < f } // MATCH /expression always evaluates to false/

func foo3(f float64) bool { return f <= f }

func foo4(f float64) bool { return f >= f }

func foo5(f float64) bool { return f == f }

func foo6(f float64) bool { return fmt.Sprintf("%s", buf1.Bytes()) == fmt.Sprintf("%s", buf1.Bytes()) }

func foo7(f float64) bool {
	return fFoo(fBar(isNaN(10.), bpar), 10000) || fFoo(fBar(isNaN(10.), bpar), 10000)
}

func foo8(f float64) bool {
	return fFoo(fBar(isNaN(10.), bpar), 10000) && fFoo(fBar(isNaN(10.), bpar), 10000)
}

func sameInt(i int) bool { return i == i } // MATCH /expression always evaluates to true/

func sameIntLess(i int) bool { return i < i } // MATCH /expression always evaluates to false/

func sameBoolAnd(b bool) bool { return b && b } // MATCH /left and right hand-side sub-expressions are the same/

func samePureIntExpression(i int) bool { return i+1 == i+1 } // MATCH /expression always evaluates to true/

func sameConvertedInt(i int) bool { return int(i) == int(i) } // MATCH /expression always evaluates to true/

var sideEffectValue int

func nextSideEffectValue() int {
	sideEffectValue++
	return sideEffectValue
}

func repeatedSideEffectCalls() bool { return nextSideEffectValue() == nextSideEffectValue() }

func repeatedChannelReceives(ch <-chan int) bool { return <-ch == <-ch }
