package interpreter

import (
	"fmt"

	"github.com/Valeron93/crafting-interpreters/scanner"
	"github.com/Valeron93/crafting-interpreters/util"
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

	return nil, util.ReportErrorOnToken(name, "undefined variable: '%v'", name.Lexeme)
}

func (e *Environment) Define(name string, value any) error {

	// if _, ok := e.variables[name]; ok {
	// 	return fmt.Errorf("'%v' is already defined", name)
	// }

	e.variables[name] = value
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

func (e *Environment) GetAt(distance int, name scanner.Token) (any, error) {
	env, err := e.ancestor(distance)
	if err != nil {
		return nil, err
	}

	value, ok := env.variables[name.Lexeme]
	if !ok {
		return nil, util.ReportErrorOnToken(name, "variable '%v' not found", name.Lexeme)
	}

	return value, nil
}

func (e *Environment) AssignAt(distance int, name scanner.Token, value any) error {
	env, err := e.ancestor(distance)
	if err != nil {
		return err
	}
	env.variables[name.Lexeme] = value
	return nil
}

func (e *Environment) ancestor(distance int) (*Environment, error) {
	env := e
	for range distance {
		env = env.enclosing
	}

	return env, nil
}
