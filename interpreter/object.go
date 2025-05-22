package interpreter

import "github.com/Valeron93/crafting-interpreters/scanner"

type Object interface {
	Set(name scanner.Token, value any) error
	Get(name scanner.Token) (any, error)
}

type IndexableObject interface {
	SetKeyValue(interpreter *Interpreter, bracket scanner.Token, key any, value any) error
	GetKeyValue(interpreter *Interpreter, bracket scanner.Token, key any) (any, error)
}
