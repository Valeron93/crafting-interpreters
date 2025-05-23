// DO NOT MODIFY!!! This file is generated from stmts.ast

package ast

import "github.com/Valeron93/crafting-interpreters/scanner"

type StmtVisitor interface {
	VisitExprStmt(*ExprStmt) (any, error)
	VisitVarStmt(*VarStmt) (any, error)
	VisitIfStmt(*IfStmt) (any, error)
	VisitBlockStmt(*BlockStmt) (any, error)
	VisitWhileStmt(*WhileStmt) (any, error)
	VisitFuncDeclStmt(*FuncDeclStmt) (any, error)
	VisitReturnStmt(*ReturnStmt) (any, error)
	VisitClassDeclStmt(*ClassDeclStmt) (any, error)
	VisitMethodDeclStmt(*MethodDeclStmt) (any, error)
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

type BlockStmt struct {
	Statements []Stmt
}

func (b *BlockStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitBlockStmt(b)
}

type WhileStmt struct {
	Condition Expr
	Body Stmt
}

func (w *WhileStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitWhileStmt(w)
}

type FuncDeclStmt struct {
	Name scanner.Token
	Params []scanner.Token
	Body []Stmt
}

func (f *FuncDeclStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitFuncDeclStmt(f)
}

type ReturnStmt struct {
	scanner.Token
	Value Expr
}

func (r *ReturnStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitReturnStmt(r)
}

type ClassDeclStmt struct {
	Name scanner.Token
	Methods []*MethodDeclStmt
	Superclass *VarExpr
}

func (c *ClassDeclStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitClassDeclStmt(c)
}

type MethodDeclStmt struct {
	Func *FuncDeclStmt
	Static bool
}

func (m *MethodDeclStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitMethodDeclStmt(m)
}

