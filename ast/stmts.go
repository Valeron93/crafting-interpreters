// DO NOT MODIFY!!! This file is generated from stmts.ast

package ast

import "github.com/Valeron93/crafting-interpreters/scanner"

type StmtVisitor interface {
	VisitExprStmt(*ExprStmt) (any, error)
	VisitPrintStmt(*PrintStmt) (any, error)
	VisitVarStmt(*VarStmt) (any, error)
	VisitIfStmt(*IfStmt) (any, error)
}

type Stmt interface {
	Accept(StmtVisitor) (any, error)
}

type ExprStmt struct {
	Expr
}

func (e *ExprStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitExprStmt(e)
}

type PrintStmt struct {
	Expr
}

func (p *PrintStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitPrintStmt(p)
}

type VarStmt struct {
	Name scanner.Token
	Init Expr
}

func (v *VarStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitVarStmt(v)
}

type IfStmt struct {
	Condition Expr
	Then Stmt
	Else Stmt
}

func (i *IfStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitIfStmt(i)
}

