package interpreter

import (
	"fmt"

	"github.com/Valeron93/crafting-interpreters/scanner"
)

type Environment struct {
	variables map[string]any
	enclosing *Environment
}

func NewEnvironment() *Environment {
	env := NewSubEnvironment(nil)
	return env
}

func NewSubEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		variables: make(map[string]any),
		enclosing: enclosing,
	}
}

func (e *Environment) Get(name scanner.Token) (any, error) {
	if value, ok := e.variables[name.Lexeme]; ok {
		return value, nil
	}

	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}

	return nil, fmt.Errorf("unknown identifier: '%v'", name.Lexeme)
}

func (e *Environment) Define(name scanner.Token, value any) error {

	if _, ok := e.variables[name.Lexeme]; ok {
		return fmt.Errorf("'%v' is already defined", name.Lexeme)
	}

	e.variables[name.Lexeme] = value
	return nil
}

func (e *Environment) Assign(target scanner.Token, value any) error {
	if _, ok := e.variables[target.Lexeme]; ok {
		e.variables[target.Lexeme] = value
		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.Assign(target, value)
	}

	return fmt.Errorf("'%v' is not defined", target.Lexeme)
}
