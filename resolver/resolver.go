package resolver

import (
	"github.com/Valeron93/crafting-interpreters/ast"
	"github.com/Valeron93/crafting-interpreters/interpreter"
	"github.com/Valeron93/crafting-interpreters/scanner"
	"github.com/Valeron93/crafting-interpreters/stack"
	"github.com/Valeron93/crafting-interpreters/util"
)

type scopeMap map[string]bool

type funcType int

const (
	functionNone funcType = iota
	functionFunc
)

type Resolver struct {
	interpreter     *interpreter.Interpreter
	scopes          stack.Stack[scopeMap]
	currentFunction funcType
	errs            []error
}

func New(i *interpreter.Interpreter) *Resolver {
	r := &Resolver{
		interpreter: i,
		errs:        make([]error, 0),
	}
	r.beginScope()
	return r
}

func (r *Resolver) ResolveStatements(stmts []ast.Stmt) []error {
	for _, stmt := range stmts {
		r.resolveStmt(stmt)
	}

	return r.errs
}

func (r *Resolver) resolveStmt(stmt ast.Stmt) {
	stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr ast.Expr) {
	expr.Accept(r)
}

func (r *Resolver) resolveFunction(params []scanner.Token, body []ast.Stmt, typ funcType) {
	enclosingFunction := r.currentFunction
	r.currentFunction = typ
	r.beginScope()
	for _, param := range params {
		r.declare(param)
		r.define(param)
	}
	r.ResolveStatements(body)
	r.endScope()
	r.currentFunction = enclosingFunction
}

func (r *Resolver) beginScope() {
	scope := make(scopeMap)
	r.scopes.Push(scope)
}

func (r *Resolver) endScope() {
	r.scopes.Pop()
}

func (r *Resolver) declare(name scanner.Token) {
	if r.scopes.Count() == 0 {
		return
	}
	scope := r.scopes.MustPeek()
	if _, ok := scope[name.Lexeme]; ok {
		r.addError(util.ReportErrorOnToken(name, "'%v' was already defined in this scope", name.Lexeme))
		return
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

func (r *Resolver) addError(err error) {
	r.errs = append(r.errs, err)
}

func (r *Resolver) ClearErrors() {
	r.errs = make([]error, 0)
}
