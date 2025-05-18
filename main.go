package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Valeron93/crafting-interpreters/scanner"
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

func runString(code string) {
	scanner := scanner.NewScanner(code)

	tokens, errors := scanner.ScanTokens()

	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	}

	for _, token := range tokens {
		fmt.Println(token)
	}
}

func runFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file: %v\n", err)
		os.Exit(1)
	}
	runString(string(bytes))
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		runString(line)
	}
}
