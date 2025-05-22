package interpreter

type Callable interface {
	Call(i *Interpreter, args []any) (any, error)
	Arity() (int, bool)
	Bind(this any) Callable
}
