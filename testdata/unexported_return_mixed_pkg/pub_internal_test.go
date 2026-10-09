package pub

import "testing"

func TestInternal(t *testing.T) {
	_ = New[int](1)
	_ = NewExported[string]("test")
	_ = ReturnParam[float64](3.14)
	_ = NewExportedWithUnexportedParam()
}
