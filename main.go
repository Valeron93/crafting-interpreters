package main

import (
	"fmt"
	"os"

	"github.com/Valeron93/crafting-interpreters/interpreter"
	"github.com/Valeron93/crafting-interpreters/parser"
	"github.com/Valeron93/crafting-interpreters/resolver"
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

func runString(i *interpreter.Interpreter, r *resolver.Resolver, code string, fileName string) bool {
	scanner := scanner.NewScanner(code)

	tokens, errs := scanner.ScanTokens()

	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Fprintf(os.Stderr, "%v:%v\n", fileName, err)
		}
		return false
	}

	p := parser.NewParser(tokens)
	expression, reporter := p.Parse()
	if reporter.HasErrors() {
		for _, err := range reporter.Errors() {
			fmt.Fprintf(os.Stderr, "%v:%v\n", fileName, err)
		}
		return false
	}

	r.ResolveStatements(expression)
	reporter = r.ErrorReporter()
	if reporter.HasErrors() {
		for _, err := range reporter.Errors() {
			fmt.Fprintf(os.Stderr, "%v:%v\n", fileName, err)
		}
		return false
	}

	if err := i.Interpret(expression); err != nil {
		fmt.Fprintf(os.Stderr, "%v:%v\n", fileName, err)
		return false
	}
	return true
}

func runFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file: %v\n", err)
		os.Exit(1)
	}
	interpreter := interpreter.New()
	resolver := resolver.New(&interpreter)
	if ok := runString(&interpreter, resolver, string(bytes), path); !ok {
		os.Exit(1)
	}
}

var replInterpreter = interpreter.New()
var replResolver = resolver.New(&replInterpreter)

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
		runString(&replInterpreter, replResolver, line, "<repl>")
		replResolver.ErrorReporter().Clear()
	}
}
