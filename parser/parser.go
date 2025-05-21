package parser

import (
	"github.com/Valeron93/crafting-interpreters/ast"
	"github.com/Valeron93/crafting-interpreters/scanner"
	"github.com/Valeron93/crafting-interpreters/util"
)

type Parser struct {
	tokens  []scanner.Token
	current int
}

func NewParser(tokens []scanner.Token) Parser {
	return Parser{
		tokens:  tokens,
		current: 1,
	}
}

func (p *Parser) match(types ...scanner.TokenType) bool {
	for _, typ := range types {
		if p.check(typ) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(typ scanner.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == typ
}

func (p *Parser) checkNext(typ scanner.TokenType) bool {

	if p.current+1 >= len(p.tokens) {
		return false
	}

	return p.tokens[p.current+1].Type == typ
}

func (p *Parser) advance() scanner.Token {
	if !p.isAtEnd() {
		p.current++
	}

	return p.prev()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == scanner.EOF
}

func (p *Parser) peek() scanner.Token {
	return p.tokens[p.current]
}

func (p *Parser) prev() scanner.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(typ scanner.TokenType, msg string) (scanner.Token, error) {
	if p.check(typ) {
		return p.advance(), nil
	} else {
		return scanner.Token{}, util.ReportErrorOnToken(p.prev(), "%v", msg)
	}
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (ast.Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(scanner.Equal) {
		equals := p.prev()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if varExpr, ok := expr.(*ast.VarExpr); ok {
			name := varExpr.Name
			return &ast.AssignExpr{
				Name:  name,
				Value: value,
			}, nil
		}
		return nil, util.ReportErrorOnToken(equals, "invalid assignment")
	}
	return expr, nil
}

func (p *Parser) or() (ast.Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}
	for p.match(scanner.Or) {
		operator := p.prev()
		right, err := p.and()
		if err != nil {
			return nil, err
		}

		expr = &ast.LogicalExpr{
			Left:     expr,
			Right:    right,
			Operator: operator,
		}
	}

	return expr, nil
}

func (p *Parser) and() (ast.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.And) {
		operator := p.prev()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}

		expr = &ast.LogicalExpr{
			Left:     expr,
			Right:    right,
			Operator: operator,
		}
	}
	return expr, nil
}

func (p *Parser) equality() (ast.Expr, error) {
	expression, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.BangEqual, scanner.EqualEqual) {
		op := p.prev()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expression = &ast.BinaryExpr{
			Left:     expression,
			Operator: op,
			Right:    right,
		}
	}

	return expression, nil
}

func (p *Parser) comparison() (ast.Expr, error) {
	expression, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.Greater, scanner.GreaterEqual, scanner.Less, scanner.LessEqual) {
		op := p.prev()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expression = &ast.BinaryExpr{
			Left:     expression,
			Operator: op,
			Right:    right,
		}
	}
	return expression, nil
}

func (p *Parser) term() (ast.Expr, error) {
	expression, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.Plus, scanner.Minus) {
		op := p.prev()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expression = &ast.BinaryExpr{
			Left:     expression,
			Operator: op,
			Right:    right,
		}
	}
	return expression, nil
}

func (p *Parser) factor() (ast.Expr, error) {
	expression, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.Star, scanner.Slash) {
		op := p.prev()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expression = &ast.BinaryExpr{
			Left:     expression,
			Operator: op,
			Right:    right,
		}
	}
	return expression, nil
}

func (p *Parser) unary() (ast.Expr, error) {
	if p.match(scanner.Bang, scanner.Minus) {
		op := p.prev()
		right, err := p.unary()
		return &ast.UnaryExpr{
			Operator: op,
			Right:    right,
		}, err
	}
	return p.call()
}

func (p *Parser) call() (ast.Expr, error) {

	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(scanner.LeftParen) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) finishCall(callee ast.Expr) (ast.Expr, error) {
	args := make([]ast.Expr, 0, 1)
	if !p.check(scanner.RightParen) {

		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		args = append(args, expr)

		for p.match(scanner.Comma) {
			if len(args) >= 127 {
				return nil, util.ReportErrorOnToken(p.peek(), "function call has a limit of 127 arguments")
			}

			expr, err := p.expression()
			if err != nil {
				return nil, err
			}
			args = append(args, expr)
		}
	}

	paren, err := p.consume(scanner.RightParen, "expected ')' after function call argument list")
	if err != nil {
		return nil, err
	}
	return &ast.CallExpr{
		Callee: callee,
		Paren:  paren,
		Args:   args,
	}, nil
}

func (p *Parser) primary() (ast.Expr, error) {
	if p.match(scanner.False) {
		return &ast.LiteralExpr{
			Value: false,
		}, nil
	}

	if p.match(scanner.True) {
		return &ast.LiteralExpr{
			Value: true,
		}, nil
	}

	if p.match(scanner.Nil) {
		return &ast.LiteralExpr{
			Value: nil,
		}, nil
	}

	if p.match(scanner.Number, scanner.String) {
		return &ast.LiteralExpr{
			Value: p.prev().Literal,
		}, nil
	}

	if p.match(scanner.Ident) {
		return &ast.VarExpr{
			Name: p.prev(),
		}, nil
	}

	if p.match(scanner.LeftParen) {
		expression, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(scanner.RightParen, "expected ')' after expression")
		if err != nil {
			return nil, err
		}
		return &ast.GroupingExpr{
			Expression: expression,
		}, nil
	}

	if p.match(scanner.Func) {
		return p.lambdaFunction()
	}

	return nil, util.ReportErrorOnToken(p.prev(), "expected expression, got '%v'", p.peek().Lexeme)
}

func (p *Parser) statement() (ast.Stmt, error) {
	if p.match(scanner.If) {
		return p.ifStatement()
	}

	if p.match(scanner.While) {
		return p.whileStatement()
	}

	if p.match(scanner.Return) {
		return p.returnStatement()
	}

	if p.match(scanner.For) {
		return p.forStatement()
	}

	if p.match(scanner.LeftBrace) {
		stmts, err := p.block()
		if err != nil {
			return nil, err
		}
		return &ast.BlockStmt{
			Statements: stmts,
		}, nil
	}

	return p.expressionStatement()
}

func (p *Parser) returnStatement() (ast.Stmt, error) {
	token := p.prev()
	var value ast.Expr

	if !p.check(scanner.Semicolon) {
		var err error
		value, err = p.expression()
		if err != nil {
			return nil, util.ReportErrorOnToken(token, "expected expression or ';' after return statement")
		}
	}

	_, err := p.consume(scanner.Semicolon, "expected ';' after return statement")
	if err != nil {
		return nil, err
	}
	return &ast.ReturnStmt{
		Token: token,
		Value: value,
	}, nil
}

func (p *Parser) forStatement() (ast.Stmt, error) {
	_, err := p.consume(scanner.LeftParen, "expected '(' after for")
	if err != nil {
		return nil, err
	}

	var init ast.Stmt

	if p.match(scanner.Semicolon) {
		init = nil
	} else if p.match(scanner.Var) {
		init, err = p.varDeclaration()
	} else {
		init, err = p.expressionStatement()
	}

	if err != nil {
		return nil, err
	}

	var cond ast.Expr
	if !p.check(scanner.Semicolon) {
		cond, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(scanner.Semicolon, "expected ';' after loop condition")
	if err != nil {
		return nil, err
	}

	var incr ast.Expr
	if !p.check(scanner.RightParen) {
		incr, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(scanner.RightParen, "expected ')' after for")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if incr != nil {
		body = &ast.BlockStmt{
			Statements: []ast.Stmt{body, &ast.ExprStmt{Expr: incr}},
		}
	}

	if cond == nil {
		cond = &ast.LiteralExpr{Value: true}
	}
	body = &ast.WhileStmt{
		Condition: cond,
		Body:      body,
	}

	if init != nil {
		body = &ast.BlockStmt{
			Statements: []ast.Stmt{
				init, body,
			},
		}
	}

	return body, nil
}

func (p *Parser) whileStatement() (ast.Stmt, error) {

	_, err := p.consume(scanner.LeftParen, "expected '(' after while")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(scanner.RightParen, "expected ')' after while")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &ast.WhileStmt{
		Condition: condition,
		Body:      body,
	}, nil
}

func (p *Parser) block() ([]ast.Stmt, error) {
	const msg = "expected '}' after block"

	stmts := make([]ast.Stmt, 0, 10)
	for !p.check(scanner.RightBrace) {
		declaration, err := p.declaration()
		if err != nil {
			return nil, util.ReportErrorOnToken(p.peek(), msg)
		}
		stmts = append(stmts, declaration)
	}
	_, err := p.consume(scanner.RightBrace, msg)
	return stmts, err

}

func (p *Parser) ifStatement() (ast.Stmt, error) {
	_, err := p.consume(scanner.LeftParen, "expected '(' after 'if'")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(scanner.RightParen, "expected ')' after if condition")
	if err != nil {
		return nil, err
	}

	thenStmt, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseStmt ast.Stmt
	if p.match(scanner.Else) {
		elseStmt, err = p.statement()
		if err != nil {
			return nil, err
		}
	}
	return &ast.IfStmt{
		Condition: condition,
		Then:      thenStmt,
		Else:      elseStmt,
	}, nil

}

func (p *Parser) declaration() (ast.Stmt, error) {
	if p.check(scanner.Func) && p.checkNext(scanner.Ident) {
		p.match(scanner.Func)
		return p.function("function")
	}

	if p.match(scanner.Class) {
		return p.classDeclaration()
	}

	if p.match(scanner.Var) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) classDeclaration() (ast.Stmt, error) {
	name, err := p.consume(scanner.Ident, "expected class name")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(scanner.LeftBrace, "expected '{' before class body")
	if err != nil {
		return nil, err
	}

	methods := make([]*ast.FuncDeclStmt, 0)
	for !p.check(scanner.RightBrace) {
		if p.match(scanner.Func) {
			fn, err := p.function("function")
			if err != nil {
				return nil, util.ReportErrorOnToken(p.prev(), "expected a function in class body")
			}
			methods = append(methods, fn)
		}
	}

	_, err = p.consume(scanner.RightBrace, "expected '}' after class body")
	return &ast.ClassDeclStmt{
		Name:    name,
		Methods: methods,
	}, nil
}

func (p *Parser) params() ([]scanner.Token, error) {

	_, err := p.consume(scanner.LeftParen, "expected '(' before parameter list")
	if err != nil {
		return nil, err
	}

	params := make([]scanner.Token, 0)
	if !p.check(scanner.RightParen) {

		param, err := p.consume(scanner.Ident, "expected parameter name")
		if err != nil {
			return nil, err
		}
		params = append(params, param)

		for p.match(scanner.Comma) {
			param, err = p.consume(scanner.Ident, "expected parameter name")
			if err != nil {
				return nil, err
			}
			params = append(params, param)
		}
	}

	_, err = p.consume(scanner.RightParen, "expected ')' after parameter list")
	return params, err
}

func (p *Parser) function(kind string) (*ast.FuncDeclStmt, error) {
	name, err := p.consume(scanner.Ident, "expected "+kind+" name")
	if err != nil {
		return nil, err
	}

	params, err := p.params()
	if err != nil {
		return nil, err
	}

	body, err := p.functionBody(p.peek(), kind)
	if err != nil {
		return nil, err
	}

	if p.prev().Type != scanner.RightBrace {
		_, err := p.consume(scanner.Semicolon, "expected ';' after function expression")
		if err != nil {
			return nil, err
		}
	}

	return &ast.FuncDeclStmt{
		Name:   name,
		Params: params,
		Body:   body,
	}, nil

}

func (p *Parser) lambdaFunction() (ast.Expr, error) {
	params, err := p.params()
	if err != nil {
		return nil, err
	}

	body, err := p.functionBody(p.peek(), "lambda")

	return &ast.LambdaExpr{
		Params: params,
		Body:   body,
	}, nil
}

func (p *Parser) functionBody(name scanner.Token, kind string) ([]ast.Stmt, error) {
	if p.match(scanner.Arrow) {
		returnExpr, err := p.expression()
		if err != nil {
			return nil, err
		}

		stmt := &ast.ReturnStmt{
			Token: name,
			Value: returnExpr,
		}

		return []ast.Stmt{stmt}, nil
	}

	_, err := p.consume(scanner.LeftBrace, "expected '{' before "+kind+" body")
	if err != nil {
		return nil, err
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (p *Parser) varDeclaration() (ast.Stmt, error) {
	name, err := p.consume(scanner.Ident, "expected variable name")
	if err != nil {
		return nil, err
	}

	var init ast.Expr

	if p.match(scanner.Equal) {
		init, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(scanner.Semicolon, "expected ';' after variable declaration")
	if err != nil {
		return nil, err
	}
	return &ast.VarStmt{
		Name: name,
		Init: init,
	}, nil
}

func (p *Parser) expressionStatement() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(scanner.Semicolon, "expected ';' after expression")
	if err != nil {
		return nil, err
	}
	return &ast.ExprStmt{
		Expr: expr,
	}, nil
}

func (p *Parser) sync() {
	p.advance()

	for !p.isAtEnd() {
		if p.prev().Type == scanner.Semicolon {
			return
		}

		switch p.peek().Type {
		case scanner.Class, scanner.Func, scanner.Var, scanner.For,
			scanner.If, scanner.While, scanner.Return:
			return
		}

		p.advance()
	}
}
func (p *Parser) Parse() ([]ast.Stmt, []error) {

	stmts := make([]ast.Stmt, 0, 100)
	errs := make([]error, 0)
	for !p.isAtEnd() {

		stmt, err := p.declaration()
		if err != nil {
			p.sync()
			errs = append(errs, err)
		} else {
			stmts = append(stmts, stmt)
		}
	}

	return stmts, errs
}
