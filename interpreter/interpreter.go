package interpreter

import (
	"errors"
	"fmt"

	"github.com/Valeron93/crafting-interpreters/ast"
	"github.com/Valeron93/crafting-interpreters/scanner"
)

type Interpreter struct {
	env *Environment
}

func (i *Interpreter) VisitAssignExpr(expr *ast.AssignExpr) (any, error) {
	value, err := i.Eval(expr.Value)
	if err != nil {
		return nil, err
	}
	err = i.env.Assign(expr.Name, value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func New() Interpreter {
	return Interpreter{
		env: NewEnvironment(),
	}
}

func (i *Interpreter) VisitVarStmt(stmt *ast.VarStmt) (any, error) {
	var value any
	var err error
	if stmt.Init != nil {
		value, err = i.Eval(stmt.Init)
		if err != nil {
			return nil, err
		}
	}
	i.env.Define(stmt.Name, value)
	return nil, nil
}

func (i *Interpreter) VisitVarExpr(expr *ast.VarExpr) (any, error) {

	value, err := i.env.Get(expr.Name)
	if err != nil {
		return nil, err
	}
	return value, nil
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

func (i *Interpreter) VisitBinaryExpr(expr *ast.BinaryExpr) (any, error) {
	left, err := i.Eval(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.Eval(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case scanner.EqualEqual:
		return left == right, nil
	case scanner.BangEqual:
		return left != right, nil
	}
	lhsFloat, lhsIsFloat := left.(float64)
	rhsFloat, rhsIsFloat := right.(float64)

	lhsString, lhsIsString := left.(string)
	rhsString, rhsIsString := right.(string)

	if lhsIsFloat && rhsIsFloat {

		switch expr.Operator.Type {
		case scanner.Plus, scanner.Minus, scanner.Star, scanner.Slash:
			return floatOperator(expr.Operator, lhsFloat, rhsFloat), nil

		case scanner.Greater, scanner.GreaterEqual, scanner.Less, scanner.LessEqual:
			return floatLogicOperator(expr.Operator, lhsFloat, rhsFloat), nil
		}
	} else if (lhsIsString && rhsIsString) && expr.Operator.Type == scanner.Plus {
		return lhsString + rhsString, nil
	} else {
		return nil, fmt.Errorf("operands '%v' and '%v' are not compatible with binary operator '%v'", left, right, expr.Operator.Lexeme)
	}

	return nil, errors.ErrUnsupported
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.UnaryExpr) (any, error) {
	right, err := i.Eval(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case scanner.Bang:
		return !(i.isTrue(right)), nil
	case scanner.Minus:
		return -(right.(float64)), nil
	}
	return nil, nil
}

func (i *Interpreter) VisitLiteralExpr(expr *ast.LiteralExpr) (any, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitGroupingExpr(expr *ast.GroupingExpr) (any, error) {
	return i.Eval(expr.Expression)
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

func (i *Interpreter) VisitExprStmt(stmt *ast.ExprStmt) (any, error) {
	_, err := i.Eval(stmt.Expr)
	return nil, err
}

func (i *Interpreter) VisitPrintStmt(stmt *ast.PrintStmt) (any, error) {
	value, err := i.Eval(stmt.Expr)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v\n", value)
	return nil, nil
}

func (i *Interpreter) VisitIfStmt(stmt *ast.IfStmt) (any, error) {

	cond, err := i.Eval(stmt.Condition)
	if err != nil {
		return nil, err
	}

	if i.isTrue(cond) {
		i.execute(stmt.Then)
	} else if stmt.Else != nil {
		i.execute(stmt.Else)
	}
	return nil, nil
}

func (i *Interpreter) VisitBlockStmt(stmt *ast.BlockStmt) (any, error) {
	subEnv := NewSubEnvironment(i.env)
	err := i.executeBlock(stmt.Statements, subEnv)
	if err != nil {
		return nil, err
	}
	return nil, nil
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

func (i *Interpreter) VisitLogicalExpr(expr *ast.LogicalExpr) (any, error) {
	left, err := i.Eval(expr.Left)
	if err != nil {
		return nil, err
	}
	if expr.Operator.Type == scanner.Or {
		if i.isTrue(left) {
			return left, nil
		}
	} else {
		if !i.isTrue(left) {
			return left, nil
		}
	}
	return i.Eval(expr.Right)
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
