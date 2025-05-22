package interpreter

import (
	"fmt"

	"github.com/Valeron93/crafting-interpreters/scanner"
	"github.com/Valeron93/crafting-interpreters/util"
)

type Class struct {
	Name        string
	Methods     map[string]ClassMethod
	Constructor Callable
}

type ClassMethod struct {
	Callable Callable
	Static   bool
}

func (c *Class) Call(i *Interpreter, args []any) (any, error) {
	instance := &ClassInstance{
		Class:  c,
		Fields: make(map[string]any),
	}

	if c.Constructor != nil {
		_, err := c.Constructor.Bind(instance).Call(i, args)
		if err != nil {
			return nil, err
		}
	}

	return instance, nil
}

func (c *Class) Arity() (int, bool) {
	if c.Constructor != nil {
		return c.Constructor.Arity()
	}
	return 0, false
}

func (c *Class) String() string {
	return fmt.Sprintf("<class %v>", c.Name)
}

func (c *Class) Bind(this any) Callable {
	return c
}

func (c *Class) FindMethod(name string) (Callable, bool) {
	if method, ok := c.Methods[name]; ok {
		return method.Callable, !method.Static
	}
	return nil, false
}

func (c *Class) Set(name scanner.Token, value any) error {
	return util.ReportErrorOnToken(name, "assigning to class object is not supported")
}

func (c *Class) Get(name scanner.Token) (any, error) {
	if method, ok := c.Methods[name.Lexeme]; ok {
		return method.Callable, nil
	}
	return nil, util.ReportErrorOnToken(name, "class '%v' has no static method '%v'", c.Name, name.Lexeme)
}
