package interpreter

import "time"

type Clock struct {
}

func (c *Clock) Call(i *Interpreter, args []any) (any, error) {
	return float64(time.Now().UnixNano()) / 1e9, nil
}

func (c *Clock) Arity() int {
	return 0
}
