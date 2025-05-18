// DO NOT MODIFY!!! This file is generated from exprs.ast

package ast

import "github.com/Valeron93/crafting-interpreters/scanner"

type ExprVisitor interface {
	VisitBinaryExpr(*BinaryExpr) any
	VisitGroupingExpr(*GroupingExpr) any
	VisitLiteralExpr(*LiteralExpr) any
	VisitUnaryExpr(*UnaryExpr) any
	VisitVarExpr(*VarExpr) any
	VisitAssignExpr(*AssignExpr) any
}

type Expr interface {
	Accept(ExprVisitor) any
}

type BinaryExpr struct {
	Right Expr
	Left Expr
	Operator scanner.Token
}

func (b *BinaryExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitBinaryExpr(b)
}

type GroupingExpr struct {
	Expression Expr
}

func (g *GroupingExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitGroupingExpr(g)
}

type LiteralExpr struct {
	Value any
}

func (l *LiteralExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitLiteralExpr(l)
}

type UnaryExpr struct {
	Operator scanner.Token
	Right Expr
}

func (u *UnaryExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitUnaryExpr(u)
}

type VarExpr struct {
	Name scanner.Token
}

func (v *VarExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitVarExpr(v)
}

type AssignExpr struct {
	Name scanner.Token
	Value Expr
}

func (a *AssignExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitAssignExpr(a)
}

