package parser

// TODO: change panics to go's errors

import (
	"fmt"

	"github.com/Valeron93/crafting-interpreters/ast"
	"github.com/Valeron93/crafting-interpreters/scanner"
)

type Parser struct {
	tokens  []scanner.Token
	current int
}

func NewParser(tokens []scanner.Token) Parser {
	return Parser{
		tokens:  tokens,
		current: 0,
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

func (p *Parser) panic(token scanner.Token, err string) {
	panic(fmt.Sprintf("error on line %v: %v", p.peek().Line, err))
}

func (p *Parser) consume(typ scanner.TokenType, err string) scanner.Token {
	if p.check(typ) {
		return p.advance()
	}

	p.panic(p.peek(), err)
	panic("unreachable")
}

func (p *Parser) expression() ast.Expr {
	return p.equality()
}

func (p *Parser) equality() ast.Expr {
	expression := p.comparison()

	for p.match(scanner.BangEqual, scanner.EqualEqual) {
		op := p.prev()
		right := p.comparison()
		expression = &ast.BinaryExpr{
			Left:     expression,
			Operator: op,
			Right:    right,
		}
	}

	return expression
}

func (p *Parser) comparison() ast.Expr {
	expression := p.term()

	for p.match(scanner.Greater, scanner.GreaterEqual, scanner.Less, scanner.LessEqual) {
		op := p.prev()
		right := p.term()
		expression = &ast.BinaryExpr{
			Left:     expression,
			Operator: op,
			Right:    right,
		}
	}
	return expression
}

func (p *Parser) term() ast.Expr {
	expression := p.factor()

	for p.match(scanner.Plus, scanner.Minus) {
		op := p.prev()
		right := p.factor()
		expression = &ast.BinaryExpr{
			Left:     expression,
			Operator: op,
			Right:    right,
		}
	}
	return expression
}

func (p *Parser) factor() ast.Expr {
	expression := p.unary()

	for p.match(scanner.Star, scanner.Slash) {
		op := p.prev()
		right := p.unary()
		expression = &ast.BinaryExpr{
			Left:     expression,
			Operator: op,
			Right:    right,
		}
	}
	return expression
}

func (p *Parser) unary() ast.Expr {
	if p.match(scanner.Bang, scanner.Minus) {
		op := p.prev()
		right := p.unary()
		return &ast.UnaryExpr{
			Operator: op,
			Right:    right,
		}
	}

	return p.primary()
}

func (p *Parser) primary() ast.Expr {
	if p.match(scanner.False) {
		return &ast.LiteralExpr{
			Value: false,
		}
	}

	if p.match(scanner.True) {
		return &ast.LiteralExpr{
			Value: true,
		}
	}

	if p.match(scanner.Nil) {
		return &ast.LiteralExpr{
			Value: nil,
		}
	}

	if p.match(scanner.Number, scanner.String) {
		return &ast.LiteralExpr{
			Value: p.prev().Literal,
		}
	}

	if p.match(scanner.LeftParen) {
		expression := p.expression()
		p.consume(scanner.RightParen, "expected ')' after expression")
		return &ast.GroupingExpr{
			Expression: expression,
		}
	}

	p.panic(p.peek(), "expected expression")
	panic("unreachable")
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
func (p *Parser) Parse() ast.Expr {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	return p.expression()
}
