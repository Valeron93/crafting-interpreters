package interpreter

import (
	"fmt"

	"github.com/Valeron93/crafting-interpreters/ast"
)

type Callable interface {
	Call(i *Interpreter, args []any) (any, error)
	Arity() (int, bool)
}

type CallableObject struct {
	Declaration *ast.FuncDeclStmt
}

func (c *CallableObject) Call(i *Interpreter, args []any) (any, error) {
	env := NewSubEnvironment(i.env)

	for i, arg := range args {
		env.Define(c.Declaration.Params[i].Lexeme, arg)
	}
	i.executeBlock(c.Declaration.Body, env)
	return nil, nil
}

func (c *CallableObject) Arity() (int, bool) {
	return len(c.Declaration.Params), false
}

func (c *CallableObject) String() string {
	return fmt.Sprintf("<fn %v>", c.Declaration.Name.Lexeme)
}
