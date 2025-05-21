package resolver

import (
	"github.com/Valeron93/crafting-interpreters/ast"
	"github.com/Valeron93/crafting-interpreters/util"
)

func (r *Resolver) VisitExprStmt(stmt *ast.ExprStmt) (any, error) {
	r.resolveExpr(stmt.Expr)
	return nil, nil
}

func (r *Resolver) VisitFuncDeclStmt(stmt *ast.FuncDeclStmt) (any, error) {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.resolveFunction(stmt.Params, stmt.Body, functionFunc)
	return nil, nil
}

func (r *Resolver) VisitIfStmt(stmt *ast.IfStmt) (any, error) {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Then)
	if stmt.Else != nil {
		r.resolveStmt(stmt.Else)
	}
	return nil, nil
}

func (r *Resolver) VisitBlockStmt(stmt *ast.BlockStmt) (any, error) {
	r.beginScope()
	r.ResolveStatements(stmt.Statements)

	r.endScope()
	return nil, nil
}

func (r *Resolver) VisitReturnStmt(stmt *ast.ReturnStmt) (any, error) {

	if r.currentFunction == functionNone {
		r.addError(util.ReportErrorOnToken(stmt.Token, "return is allowed only in functions"))
		return nil, nil
	}

	if stmt.Value != nil {
		r.resolveExpr(stmt.Value)
	}
	return nil, nil
}

func (r *Resolver) VisitVarStmt(stmt *ast.VarStmt) (any, error) {
	r.declare(stmt.Name)
	if stmt.Init != nil {
		r.resolveExpr(stmt.Init)
	}
	r.define(stmt.Name)
	return nil, nil
}

func (r *Resolver) VisitWhileStmt(stmt *ast.WhileStmt) (any, error) {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Body)
	return nil, nil
}

func (r *Resolver) VisitClassDeclStmt(stmt *ast.ClassDeclStmt) (any, error) {
	r.declare(stmt.Name)
	r.define(stmt.Name)
	return nil, nil
}
