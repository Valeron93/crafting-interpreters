package resolver

import (
	"github.com/Valeron93/crafting-interpreters/ast"
	"github.com/Valeron93/crafting-interpreters/scanner"
	"github.com/Valeron93/crafting-interpreters/util"
)

func (r *Resolver) VisitAssignExpr(expr *ast.AssignExpr) (any, error) {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) VisitBinaryExpr(expr *ast.BinaryExpr) (any, error) {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil, nil
}

func (r *Resolver) VisitCallExpr(expr *ast.CallExpr) (any, error) {
	r.resolveExpr(expr.Callee)

	for _, arg := range expr.Args {
		r.resolveExpr(arg)
	}

	return nil, nil
}

func (r *Resolver) VisitGroupingExpr(expr *ast.GroupingExpr) (any, error) {
	r.resolveExpr(expr.Expression)
	return nil, nil
}

func (r *Resolver) VisitLambdaExpr(expr *ast.LambdaExpr) (any, error) {
	r.resolveFunction(expr.Params, expr.Body, functionFunc)
	return nil, nil
}

func (r *Resolver) VisitLiteralExpr(*ast.LiteralExpr) (any, error) {
	return nil, nil
}

func (r *Resolver) VisitLogicalExpr(expr *ast.LogicalExpr) (any, error) {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil, nil
}

func (r *Resolver) VisitUnaryExpr(expr *ast.UnaryExpr) (any, error) {
	r.resolveExpr(expr.Right)
	return nil, nil
}

func (r *Resolver) VisitVarExpr(expr *ast.VarExpr) (any, error) {

	if !r.scopes.Empty() {
		if value, ok := r.scopes.MustPeek()[expr.Name.Lexeme]; ok && value == false {
			r.addError(util.ReportErrorOnToken(expr.Name, "can't read local var in its own initializer"))
		}
	}

	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) VisitGetExpr(expr *ast.GetExpr) (any, error) {
	r.resolveExpr(expr.Object)
	return nil, nil
}

func (r *Resolver) VisitSetExpr(expr *ast.SetExpr) (any, error) {
	r.resolveExpr(expr.Value)
	r.resolveExpr(expr.Object)
	return nil, nil
}

func (r *Resolver) VisitThisExpr(expr *ast.ThisExpr) (any, error) {

	if r.currentClass == classNone {
		r.addError(util.ReportErrorOnToken(expr.Keyword, "cannot use this outside of class method"))
		return nil, nil
	}

	r.resolveLocal(expr, expr.Keyword)
	return nil, nil
}

func (r *Resolver) resolveLocal(expr ast.Expr, name scanner.Token) {
	for i := r.scopes.Count() - 1; i >= 0; i-- {
		scope := r.scopes.GetIdx(i)

		if _, ok := scope[name.Lexeme]; ok {
			r.interpreter.Resolve(expr, r.scopes.Count()-1-i)
			return
		}
	}

	if !r.interpreter.GlobalExists(name.Lexeme) {
		r.addError(util.ReportErrorOnToken(name, "undefined variable: %v", name.Lexeme))
	}
}

func (r *Resolver) VisitGetKeyExpr(expr *ast.GetKeyExpr) (any, error) {
	r.resolveExpr(expr.Object)
	r.resolveExpr(expr.Key)
	return nil, nil
}

func (r *Resolver) VisitSetKeyExpr(expr *ast.SetKeyExpr) (any, error) {
	r.resolveExpr(expr.Object)
	r.resolveExpr(expr.Key)
	r.resolveExpr(expr.Value)
	return nil, nil
}

func (r *Resolver) VisitSuperExpr(expr *ast.SuperExpr) (any, error) {

	if r.currentClass == classNone {
		r.addError(util.ReportErrorOnToken(expr.Keyword, "cannot use 'super' outside of class"))
	} else if r.currentClass != classSubclass {
		r.addError(util.ReportErrorOnToken(expr.Keyword, "cannot use 'super' in a class without superclass"))
	}

	r.resolveLocal(expr, expr.Keyword)
	return nil, nil
}
