package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/Valeron93/crafting-interpreters/parser"
	"github.com/Valeron93/crafting-interpreters/scanner"
	"github.com/chzyer/readline"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Fprintf(os.Stderr, "usage: %v [filename]\n", os.Args[0])
		os.Exit(64)
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runString(i *Interpreter, code string) error {
	scanner := scanner.NewScanner(code)

	tokens, errs := scanner.ScanTokens()

	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		return errors.New("failed to scan")
	}

	p := parser.NewParser(tokens)
	expression, errs := p.Parse()
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		return errors.New("failed to parse")
	}

	if expression != nil {
		if err := i.Interpret(expression); err != nil {
			return err
		}
	}
	return nil
}

func runFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file: %v\n", err)
		os.Exit(1)
	}
	interpreter := NewInterpreter()
	if err := runString(&interpreter, string(bytes)); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\n", err)
		os.Exit(1)
	}
}

var replInterpreter = NewInterpreter()

func runPrompt() {

	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		if err := runString(&replInterpreter, line); err != nil {
			fmt.Fprintf(os.Stderr, "runtime error: %v\n", err)
		}
	}
}
