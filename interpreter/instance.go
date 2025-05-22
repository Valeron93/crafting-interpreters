package interpreter

import (
	"fmt"

	"github.com/Valeron93/crafting-interpreters/scanner"
	"github.com/Valeron93/crafting-interpreters/util"
)

type ClassInstance struct {
	Class  *Class
	Fields map[string]any
}

func (c *ClassInstance) Set(name scanner.Token, value any) error {
	c.Fields[name.Lexeme] = value
	return nil
}

func (c *ClassInstance) Get(name scanner.Token) (any, error) {
	if field, ok := c.Fields[name.Lexeme]; ok {
		return field, nil
	}

	if method, ok := c.Class.FindMethod(name.Lexeme); ok {
		return method.Bind(c), nil
	}

	return nil, util.ReportErrorOnToken(name, "undefined field '%v'", name.Lexeme)
}

func (c *ClassInstance) String() string {
	return fmt.Sprintf("<%v obj %p>", c.Class.Name, c)
}

func (c *ClassInstance) SetKeyValue(interpreter *Interpreter, bracket scanner.Token, key any, value any) error {

	if method, ok := c.Class.FindMethod("__set"); ok {

		_, err := method.Bind(c).Call(interpreter, []any{key, value})

		if err != nil {
			return err
		}

		return nil
	}
	return util.ReportErrorOnToken(bracket, "class '%v' does not have 'fn __set(key, value)' method", c.Class.Name)

}

func (c *ClassInstance) GetKeyValue(interpreter *Interpreter, bracket scanner.Token, key any) (any, error) {

	if method, ok := c.Class.FindMethod("__get"); ok {
		value, err := method.Bind(c).Call(interpreter, []any{key})
		if err != nil {
			return nil, err
		}

		return value, nil
	}

	return nil, util.ReportErrorOnToken(bracket, "class '%v' does not have 'fn __get(key)' method", c.Class.Name)
}
