package pub

type impl[T any] struct{ val T }

type Exported[T any] struct{ Val T }

type hidden struct{ n int }

func New[T any](v T) *impl[T] { return &impl[T]{val: v} }

func NewExported[T any](v T) *Exported[T] { return &Exported[T]{Val: v} }

func ReturnParam[T any](v T) T { return v }

func NewExportedWithUnexportedParam() *Exported[hidden] { return &Exported[hidden]{Val: hidden{n: 1}} }
