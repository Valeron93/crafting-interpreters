package scanner

import (
	"fmt"
	"strconv"
)

type Scanner struct {
	source []rune
	tokens []Token

	start   int
	current int
	line    int
	column  int
}

func NewScanner(source string) Scanner {
	return Scanner{
		source:  []rune(source),
		tokens:  []Token{},
		start:   0,
		current: 0,
		line:    1,
		column:  1,
	}
}

func (s *Scanner) ScanTokens() ([]Token, []error) {

	errors := []error{}

	for !s.isAtEnd() {
		s.start = s.current
		err := s.scanToken()
		if err != nil {
			errors = append(errors, err)
		}

	}

	s.tokens = append(s.tokens, Token{
		Type: EOF,
		Line: s.line,
	})

	return s.tokens, errors
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) addToken(typ TokenType) {
	s.addTokenLiteral(typ, nil)
}

func (s *Scanner) addTokenLiteral(typ TokenType, literal any) {
	text := string(s.source[s.start:s.current])
	s.tokens = append(s.tokens, Token{
		Type:    typ,
		Lexeme:  text,
		Literal: literal,
		Line:    s.line,
		Column:  s.column - len(text) + 1,
	})
}

func (s *Scanner) advance() rune {
	r := s.source[s.current]
	s.current++
	s.column++
	return r
}

func (s *Scanner) scanToken() error {

	c := s.advance()

	switch c {
	case '(':
		s.addToken(LeftParen)
	case ')':
		s.addToken(RightParen)
	case '{':
		s.addToken(LeftBrace)
	case '}':
		s.addToken(RightBrace)
	case ',':
		s.addToken(Comma)
	case '.':
		s.addToken(Dot)
	case '-':
		s.addToken(Minus)
	case '+':
		s.addToken(Plus)
	case ';':
		s.addToken(Semicolon)
	case '*':
		s.addToken(Star)

	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(Slash)
		}

	case ' ', '\r', '\t':
		break

	case '\n':
		s.line++
		s.column = 0

	case '"':
		err := s.string()
		if err != nil {
			return err
		}

	case '!':
		if s.match('=') {
			s.addToken(BangEqual)
		} else {
			s.addToken(Bang)
		}

	case '=':
		if s.match('=') {
			s.addToken(EqualEqual)
		} else if s.match('>') {
			s.addToken(Arrow)
		} else {
			s.addToken(Equal)
		}

	case '<':
		if s.match('=') {
			s.addToken(LessEqual)
		} else {
			s.addToken(Less)
		}

	case '>':
		if s.match('=') {
			s.addToken(GreaterEqual)
		} else {
			s.addToken(Greater)
		}

	default:
		if isDigit(c) {
			err := s.number()
			if err != nil {
				return err
			}
		} else if isAlpha(c) {
			s.identifier()
		} else {
			return s.error(fmt.Sprintf("unexpected character '%c'", c))
		}
	}
	return nil
}

func (s *Scanner) match(expected rune) bool {

	if s.isAtEnd() {
		return false
	}

	if s.source[s.current] != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() rune {

	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

func (s *Scanner) string() error {

	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		return s.error("unterminated string")
	}
	s.advance()

	value := s.source[s.start+1 : s.current-1]
	s.addTokenLiteral(String, string(value))

	return nil
}

func (s *Scanner) error(err string) error {
	return fmt.Errorf("%v:%v: %v", s.line, s.column, err)
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r == '_')
}

func isAlphaNumeric(r rune) bool {
	return isAlpha(r) || isDigit(r)
}

func (s *Scanner) number() error {

	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}
	value, err := strconv.ParseFloat(string(s.source[s.start:s.current]), 64)
	if err != nil {
		return err
	}
	s.addTokenLiteral(Number, value)

	return nil
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	// Check if the identifier is a keyword
	text := string(s.source[s.start:s.current])
	if typ, ok := keywords[text]; ok {
		s.addToken(typ)
		return
	}

	s.addToken(Ident)
}
