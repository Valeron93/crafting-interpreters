// DO NOT MODIFY!!! This file is generated from exprs.ast

package expr

import "github.com/Valeron93/crafting-interpreters/scanner"

type ExprVisitor interface {
	VisitBinary(*Binary) any
	VisitGrouping(*Grouping) any
	VisitLiteral(*Literal) any
	VisitUnary(*Unary) any
}

type Expr interface {
	Accept(ExprVisitor) any
}

type Binary struct {
	Right Expr
	Left Expr
	Operator scanner.Token
}

func (b *Binary) Accept(visitor ExprVisitor) any {
	return visitor.VisitBinary(b)
}

type Grouping struct {
	Expression Expr
}

func (g *Grouping) Accept(visitor ExprVisitor) any {
	return visitor.VisitGrouping(g)
}

type Literal struct {
	Value any
}

func (l *Literal) Accept(visitor ExprVisitor) any {
	return visitor.VisitLiteral(l)
}

type Unary struct {
	Operator scanner.Token
	Right Expr
}

func (u *Unary) Accept(visitor ExprVisitor) any {
	return visitor.VisitUnary(u)
}

