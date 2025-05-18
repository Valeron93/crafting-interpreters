package main

import (
	"fmt"

	"github.com/Valeron93/crafting-interpreters/scanner"
)

type Environment struct {
	Globals map[string]any
}

func (e *Environment) Get(name scanner.Token) (any, error) {
	if value, ok := e.Globals[name.Lexeme]; ok {
		return value, nil
	}
	return nil, fmt.Errorf("unknown identifier: '%v'", name.Lexeme)
}

func (e *Environment) Define(name scanner.Token, value any) error {

	if _, ok := e.Globals[name.Lexeme]; ok {
		return fmt.Errorf("'%v' is already defined", name.Lexeme)
	}

	e.Globals[name.Lexeme] = value
	return nil
}

func (e *Environment) Assign(target scanner.Token, value any) error {
	if _, ok := e.Globals[target.Lexeme]; ok {
		e.Globals[target.Lexeme] = value
		return nil
	}
	return fmt.Errorf("'%v' is not defined", target.Lexeme)
}
