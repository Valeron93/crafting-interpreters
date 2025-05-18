package main

import (
	"fmt"
	"strings"

	"github.com/Valeron93/crafting-interpreters/ast"
)

type AstPrinter struct {
}

func (a *AstPrinter) VisitBinaryExpr(e *ast.BinaryExpr) any {
	return a.parenthesize(e.Operator.Lexeme, e.Left, e.Right)
}

func (a *AstPrinter) VisitGroupingExpr(e *ast.GroupingExpr) any {
	return a.parenthesize("group", e.Expression)
}

func (a *AstPrinter) VisitLiteralExpr(e *ast.LiteralExpr) any {
	return fmt.Sprintf("%v", e.Value)
}

func (a *AstPrinter) VisitUnaryExpr(e *ast.UnaryExpr) any {
	return a.parenthesize(e.Operator.Lexeme, e.Right)
}

func (a *AstPrinter) Print(e ast.Expr) any {
	return e.Accept(a)
}

func (a *AstPrinter) parenthesize(name string, exprs ...ast.Expr) string {
	var sb strings.Builder

	sb.WriteRune('(')
	sb.WriteString(name)
	for _, expr := range exprs {
		sb.WriteRune(' ')
		sb.WriteString(expr.Accept(a).(string))
	}
	sb.WriteRune(')')
	return sb.String()
}