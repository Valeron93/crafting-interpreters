// DO NOT MODIFY!!! This file is generated from stmts.ast

package ast

import "github.com/Valeron93/crafting-interpreters/scanner"

type StmtVisitor interface {
	VisitExprStmt(*ExprStmt) any
	VisitPrintStmt(*PrintStmt) any
	VisitVarStmt(*VarStmt) any
}

type Stmt interface {
	Accept(StmtVisitor) any
}

type ExprStmt struct {
	Expr
}

func (e *ExprStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitExprStmt(e)
}

type PrintStmt struct {
	Expr
}

func (p *PrintStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitPrintStmt(p)
}

type VarStmt struct {
	Name scanner.Token
	Init Expr
}

func (v *VarStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitVarStmt(v)
}

