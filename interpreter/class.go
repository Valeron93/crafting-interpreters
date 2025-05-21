package interpreter

import "fmt"

type Class struct {
	Name    string
	Methods map[string]*CallableObject
}

func (c *Class) Call(i *Interpreter, args []any) (any, error) {
	return &ClassInstance{
		Class:  c,
		Fields: make(map[string]any),
	}, nil
}

func (c *Class) Arity() (int, bool) {
	return 0, false
}

func (c *Class) String() string {
	return fmt.Sprintf("<class %v>", c.Name)
}

func (c *Class) FindMethod(name string) (*CallableObject, bool) {
	if method, ok := c.Methods[name]; ok {
		return method, true
	}
	return nil, false
}
