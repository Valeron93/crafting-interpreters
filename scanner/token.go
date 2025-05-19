package scanner

import "fmt"

//go:generate stringer -type=TokenType
type TokenType int32

const (
	LeftParen TokenType = iota
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star

	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual

	Ident
	String
	Number

	And
	Class
	Else
	False
	Func
	For
	If
	Nil
	Or
	Return
	Super
	This
	True
	Var
	While
	EOF
)

var keywords = map[string]TokenType{
	"and":    And,
	"class":  Class,
	"else":   Else,
	"false":  False,
	"for":    For,
	"fn":     Func,
	"if":     If,
	"null":   Nil,
	"or":     Or,
	"return": Return,
	"super":  Super,
	"this":   This,
	"true":   True,
	"let":    Var,
	"while":  While,
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	Line    int
}

func (t Token) String() string {
	if t.Literal != nil {
		return fmt.Sprintf("%v {t: `%v` val: %v}", t.Type, t.Lexeme, t.Literal)
	}

	return fmt.Sprintf("%v {t: `%v`}", t.Type, t.Lexeme)
}
