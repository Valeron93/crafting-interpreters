package interpreter

import "fmt"

type Class struct {
	Name        string
	Methods     map[string]*CallableObject
	Constructor *CallableObject
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

func (c *Class) FindMethod(name string) (*CallableObject, bool) {
	if method, ok := c.Methods[name]; ok {
		return method, true
	}
	return nil, false
}
