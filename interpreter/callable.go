package interpreter

import (
	"fmt"

	"github.com/Valeron93/crafting-interpreters/ast"
	"github.com/Valeron93/crafting-interpreters/scanner"
	"github.com/Valeron93/crafting-interpreters/util"
)

type Callable interface {
	Call(i *Interpreter, args []any) (any, error)
	Arity() (int, bool)
}

type CallableObject struct {
	Declaration *ast.FuncDeclStmt
	Closure     *Environment
}

type Class struct {
	Name    string
	Methods map[string]*CallableObject
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

func (c *Class) Call(i *Interpreter, args []any) (any, error) {
	return &ClassInstance{
		Class:  c,
		Fields: make(map[string]any),
	}, nil
}

func (c *Class) Arity() (int, bool) {
	return 0, false
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

func (c *CallableObject) Bind(instance *ClassInstance) *CallableObject {
	env := NewSubEnvironment(c.Closure)

	env.Define("this", instance)
	return &CallableObject{
		Declaration: c.Declaration,
		Closure:     env,
	}
}
