package main

import (
	"fmt"
	"strings"

	"github.com/Valeron93/crafting-interpreters/expr"
)

type AstPrinter struct {
}

func (a *AstPrinter) VisitBinary(e *expr.Binary) any {
	return a.parenthesize(e.Operator.Lexeme, e.Left, e.Right)
}

func (a *AstPrinter) VisitGrouping(e *expr.Grouping) any {
	return a.parenthesize("group", e.Expression)
}

func (a *AstPrinter) VisitLiteral(e *expr.Literal) any {
	return fmt.Sprintf("%v", e.Value)
}

func (a *AstPrinter) VisitUnary(e *expr.Unary) any {
	return a.parenthesize(e.Operator.Lexeme, e.Right)
}

func (a *AstPrinter) Print(e expr.Expr) any {
	return e.Accept(a)
}

func (a *AstPrinter) parenthesize(name string, exprs ...expr.Expr) string {
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