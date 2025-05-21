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

func (c *ClassInstance) Set(name scanner.Token, value any) {
	c.Fields[name.Lexeme] = value
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
