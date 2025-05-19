package interpreter

import (
	"github.com/Valeron93/crafting-interpreters/ast"
	"github.com/Valeron93/crafting-interpreters/scanner"
)

type Interpreter struct {
	env     *Environment
	globals *Environment
}

func New() Interpreter {
	env := NewEnvironment()
	i := Interpreter{
		env:     env,
		globals: env,
	}
	i.env.Define("clock", &ClockFunction{})
	i.env.Define("print", &PrintFunction{})

	return i
}

func floatOperator(operator scanner.Token, lhs float64, rhs float64) float64 {

	switch operator.Type {
	case scanner.Minus:
		return lhs - rhs

	case scanner.Plus:
		return lhs + rhs

	case scanner.Slash:
		return lhs / rhs

	case scanner.Star:
		return lhs * rhs
	}

	panic("unreachable: floatOperator")
}

func floatLogicOperator(operator scanner.Token, lhs float64, rhs float64) bool {
	switch operator.Type {
	case scanner.Greater:
		return lhs > rhs

	case scanner.GreaterEqual:
		return lhs >= rhs

	case scanner.Less:
		return lhs < rhs

	case scanner.LessEqual:
		return lhs <= rhs

	}
	panic("unreachable: floatLogicOperator")
}

func (i *Interpreter) Eval(expr ast.Expr) (any, error) {
	return expr.Accept(i)
}

func (i *Interpreter) isTrue(obj any) bool {
	if obj == nil {
		return false
	}
	if b, ok := obj.(bool); ok {
		return b
	}
	return true
}

func (i *Interpreter) executeBlock(stmts []ast.Stmt, env *Environment) error {
	prevEnv := i.env

	defer func() {
		i.env = prevEnv
	}()

	i.env = env
	for _, stmt := range stmts {
		err := i.execute(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) execute(stmt ast.Stmt) error {
	_, err := stmt.Accept(i)
	return err
}

func (i *Interpreter) Interpret(stmts []ast.Stmt) error {

	for _, stmt := range stmts {
		if _, err := stmt.Accept(i); err != nil {
			return err
		}
	}

	return nil
}
