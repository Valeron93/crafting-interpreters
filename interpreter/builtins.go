package interpreter

import (
	"fmt"
	"time"
)

type ClockFunction struct {
}

func (c *ClockFunction) Call(i *Interpreter, args []any) (any, error) {
	return float64(time.Now().UnixNano()) / 1e9, nil
}

func (c *ClockFunction) Arity() (int, bool) {
	return 0, false
}

func (c *ClockFunction) Bind(this any) Callable {
	return c
}

type PrintFunction struct {
}

func (p *PrintFunction) Call(i *Interpreter, args []any) (any, error) {
	for _, arg := range args {
		fmt.Printf("%v", arg)
	}
	fmt.Println()
	return nil, nil
}

func (p *PrintFunction) Arity() (int, bool) {
	return 0, true
}

func (p *PrintFunction) Bind(this any) Callable {
	return p
}
