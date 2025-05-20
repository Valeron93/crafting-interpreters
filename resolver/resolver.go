package resolver

import (
	"github.com/Valeron93/crafting-interpreters/ast"
	"github.com/Valeron93/crafting-interpreters/interpreter"
	"github.com/Valeron93/crafting-interpreters/scanner"
	"github.com/Valeron93/crafting-interpreters/stack"
	"github.com/Valeron93/crafting-interpreters/util"
)

type scopeMap map[string]bool

type Resolver struct {
	interpreter *interpreter.Interpreter
	scopes      stack.Stack[scopeMap]
	errs        util.TokenErrorReporter
}

func New(i *interpreter.Interpreter) *Resolver {
	r := &Resolver{
		interpreter: i,
		errs:        util.NewTokenErrorReporter(),
	}
	return r
}

func (r *Resolver) ErrorReporter() *util.TokenErrorReporter {
	return &r.errs
}

func (r *Resolver) ResolveStatements(stmts []ast.Stmt) {
	for _, stmt := range stmts {
		r.resolveStmt(stmt)
	}
}

func (r *Resolver) resolveStmt(stmt ast.Stmt) {
	stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr ast.Expr) {
	expr.Accept(r)
}

func (r *Resolver) resolveFunction(params []scanner.Token, body []ast.Stmt) error {
	r.beginScope()
	for _, param := range params {
		r.declare(param)
		r.define(param)
	}
	r.ResolveStatements(body)
	r.endScope()
	return nil
}

func (r *Resolver) beginScope() {
	scope := make(scopeMap)
	r.scopes.Push(scope)
}

func (r *Resolver) endScope() error {
	r.scopes.Pop()
	return nil
}

func (r *Resolver) declare(name scanner.Token) {
	if r.scopes.Count() == 0 {
		return
	}
	scope := r.scopes.MustPeek()
	if _, ok := scope[name.Lexeme]; ok {
		r.errs.Report(name, "'%v' was already defined in this scope", name.Lexeme)
	}
	scope[name.Lexeme] = false
}

func (r *Resolver) define(name scanner.Token) {
	if r.scopes.Empty() {
		return
	}
	scope := r.scopes.MustPeek()
	scope[name.Lexeme] = true
}
