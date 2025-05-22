package interpreter

import (
	"fmt"

	"github.com/Valeron93/crafting-interpreters/ast"
)

type CallableObject struct {
	Declaration *ast.FuncDeclStmt
	Closure     *Environment
}

func (c *CallableObject) Bind(this any) Callable {
	env := NewSubEnvironment(c.Closure)

	env.Define("this", this)
	return &CallableObject{
		Declaration: c.Declaration,
		Closure:     env,
	}
}

func (c *CallableObject) Call(i *Interpreter, args []any) (any, error) {
	env := NewSubEnvironment(c.Closure)

	for i, arg := range args {
		env.Define(c.Declaration.Params[i].Lexeme, arg)
	}
	err := i.executeBlock(c.Declaration.Body, env)

	if err != nil {
		if functionReturn, ok := err.(*FunctionReturn); ok {
			return functionReturn.Value, nil
		} else {
			return nil, err
		}
	}

	return nil, nil
}

func (c *CallableObject) Arity() (int, bool) {
	return len(c.Declaration.Params), false
}

func (c *CallableObject) String() string {
	return fmt.Sprintf("<fn %v %p>", c.Declaration.Name.Lexeme, c)
}
