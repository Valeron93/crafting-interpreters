package main

import (
	"fmt"

	"github.com/Valeron93/crafting-interpreters/ast"
	"github.com/Valeron93/crafting-interpreters/scanner"
)

type Interpreter struct {
	env Environment
}



func (i *Interpreter) VisitAssignExpr(expr *ast.AssignExpr) any {
	value := i.Eval(expr.Value)
	i.env.Assign(expr.Name, value)
	return value
}

func NewInterpreter() Interpreter {
	return Interpreter{
		env: Environment{
			Globals: make(map[string]any),
		},
	}
}

func (i *Interpreter) VisitVarStmt(stmt *ast.VarStmt) any {
	var value any
	if stmt.Init != nil {
		value = i.Eval(stmt.Init)
	}
	i.env.Define(stmt.Name, value)
	return nil
}

func (i *Interpreter) VisitVarExpr(expr *ast.VarExpr) any {

	value, err := i.env.Get(expr.Name)
	if err != nil {
		panic(err)
	}
	return value
}

func (i *Interpreter) VisitBinaryExpr(expr *ast.BinaryExpr) any {
	left := i.Eval(expr.Left)
	right := i.Eval(expr.Right)

	switch expr.Operator.Type {

	case scanner.Minus:
		return left.(float64) - right.(float64)

	case scanner.Slash:
		return left.(float64) / right.(float64)

	case scanner.Star:
		return left.(float64) * right.(float64)

	case scanner.Plus:
		stringLeft, okLeftString := left.(string)
		stringRight, okRightString := right.(string)

		if okLeftString && okRightString {
			return stringLeft + stringRight
		}

		floatLeft, okLeftFloat := left.(float64)
		floatRight, okRightFloat := right.(float64)

		if okLeftFloat && okRightFloat {
			return floatLeft + floatRight
		}
		panic(fmt.Sprintf("binary operator + is not compatible with '%#v' and '%#v'", left, right))

	case scanner.Greater:
		return left.(float64) > right.(float64)

	case scanner.GreaterEqual:
		return left.(float64) >= right.(float64)

	case scanner.Less:
		return left.(float64) < right.(float64)
	case scanner.LessEqual:
		return left.(float64) <= right.(float64)
	case scanner.EqualEqual:
		return left == right
	case scanner.BangEqual:
		return left != right
	}
	panic("unreachable")
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.UnaryExpr) any {
	right := i.Eval(expr.Right)

	switch expr.Operator.Type {
	case scanner.Bang:
		return !(i.isTrue(right))
	case scanner.Minus:
		return -(right.(float64))
	}
	return nil
}

func (i *Interpreter) VisitLiteralExpr(expr *ast.LiteralExpr) any {
	return expr.Value
}

func (i *Interpreter) VisitGroupingExpr(expr *ast.GroupingExpr) any {
	return i.Eval(expr.Expression)
}

func (i *Interpreter) Eval(expr ast.Expr) any {
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

func (i *Interpreter) VisitExprStmt(stmt *ast.ExprStmt) any {
	i.Eval(stmt.Expr)
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *ast.PrintStmt) any {
	value := i.Eval(stmt.Expr)
	fmt.Printf("%v\n", value)
	return nil
}

func (i *Interpreter) VisitIfStmt(stmt *ast.IfStmt) any {
	
	if i.isTrue(i.Eval(stmt.Condition)) {
		i.execute(stmt.Then)
	} else if (stmt.Else != nil) {
		i.execute(stmt.Else)
	}
	return nil
}

func (i *Interpreter) execute(stmt ast.Stmt) {
	stmt.Accept(i)
}

func (i *Interpreter) Interpret(stmts []ast.Stmt) error {

	for _, stmt := range stmts {
		stmt.Accept(i)
	}

	return nil
}
