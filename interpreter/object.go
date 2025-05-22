package interpreter

import "github.com/Valeron93/crafting-interpreters/scanner"

type Object interface {
	Set(name scanner.Token, value any) error
	Get(name scanner.Token) (any, error)
}
