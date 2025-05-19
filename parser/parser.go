package parser

import (
	"fmt"

	"github.com/Valeron93/crafting-interpreters/ast"
	"github.com/Valeron93/crafting-interpreters/scanner"
)

type Parser struct {
	tokens  []scanner.Token
	current int
	errors  []error
}

func NewParser(tokens []scanner.Token) Parser {
	return Parser{
		tokens:  tokens,
		current: 0,
		errors:  make([]error, 0),
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

func (p *Parser) reportError(token scanner.Token, msg string) error {
	err := fmt.Errorf("line %v: %v", token.Line, msg)
	p.errors = append(p.errors, err)
	return err
}

// TODO: fix
func (p *Parser) consume(typ scanner.TokenType, msg string) (scanner.Token, error) {
	if p.check(typ) {
		return p.advance(), nil
	}
	if msg == "" {
		return scanner.Token{}, p.reportError(p.peek(), fmt.Sprintf("expected %v, got: %v", typ, p.peek().Type))
	} else {
		return scanner.Token{}, p.reportError(p.peek(), msg)
	}
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (ast.Expr, error) {
	expr, err := p.equality()
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
		return nil, p.reportError(equals, "invalid assignment")
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

	return p.primary()
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

	return nil, p.reportError(p.peek(), "expected expression")
}

func (p *Parser) statement() (ast.Stmt, error) {
	if p.match(scanner.If) {
		return p.ifStatement()
	}

	if p.match(scanner.Print) {
		return p.printStatement()
	}
	return p.expressionStatement()
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
		Then: thenStmt,
		Else: elseStmt,
	}, nil

}

func (p *Parser) declaration() (ast.Stmt, error) {

	if p.match(scanner.Var) {
		return p.varDeclaration()
	}
	return p.statement()
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

func (p *Parser) printStatement() (ast.Stmt, error) {

	_, err := p.consume(scanner.LeftParen, "expected '(' after print")
	if err != nil {
		return nil, err
	}

	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(scanner.RightParen, "expected ')' after expression")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(scanner.Semicolon, "expected ';' after statement")
	if err != nil {
		return nil, err
	}

	return &ast.PrintStmt{
		Expr: value,
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
			scanner.If, scanner.While, scanner.Print, scanner.Return:
			return
		}

		p.advance()
	}
}
func (p *Parser) Parse() ([]ast.Stmt, []error) {

	stmts := make([]ast.Stmt, 0, 100)

	for !p.isAtEnd() {

		stmt, err := p.declaration()
		if err != nil {
			p.sync()
		} else {
			stmts = append(stmts, stmt)
		}
	}

	return stmts, p.errors
}
