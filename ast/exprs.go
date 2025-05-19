// DO NOT MODIFY!!! This file is generated from exprs.ast

package ast

import "github.com/Valeron93/crafting-interpreters/scanner"

type ExprVisitor interface {
	VisitBinaryExpr(*BinaryExpr) (any, error)
	VisitLogicalExpr(*LogicalExpr) (any, error)
	VisitGroupingExpr(*GroupingExpr) (any, error)
	VisitLiteralExpr(*LiteralExpr) (any, error)
	VisitUnaryExpr(*UnaryExpr) (any, error)
	VisitVarExpr(*VarExpr) (any, error)
	VisitAssignExpr(*AssignExpr) (any, error)
	VisitCallExpr(*CallExpr) (any, error)
	VisitLambdaExpr(*LambdaExpr) (any, error)
}

type Expr interface {
	Accept(ExprVisitor) (any, error)
}

type BinaryExpr struct {
	Right Expr
	Left Expr
	Operator scanner.Token
}

func (b *BinaryExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitBinaryExpr(b)
}

type LogicalExpr struct {
	Right Expr
	Left Expr
	Operator scanner.Token
}

func (l *LogicalExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitLogicalExpr(l)
}

type GroupingExpr struct {
	Expression Expr
}

func (g *GroupingExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitGroupingExpr(g)
}

type LiteralExpr struct {
	Value any
}

func (l *LiteralExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitLiteralExpr(l)
}

type UnaryExpr struct {
	Operator scanner.Token
	Right Expr
}

func (u *UnaryExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitUnaryExpr(u)
}

type VarExpr struct {
	Name scanner.Token
}

func (v *VarExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitVarExpr(v)
}

type AssignExpr struct {
	Name scanner.Token
	Value Expr
}

func (a *AssignExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitAssignExpr(a)
}

type CallExpr struct {
	Callee Expr
	Paren scanner.Token
	Args []Expr
}

func (c *CallExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitCallExpr(c)
}

type LambdaExpr struct {
	Params []scanner.Token
	Body []Stmt
}

func (l *LambdaExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitLambdaExpr(l)
}

