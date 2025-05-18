package main

import (
	"fmt"

	"github.com/Valeron93/crafting-interpreters/expr"
	"github.com/Valeron93/crafting-interpreters/scanner"
)

func main() {
	expr := expr.Binary{
		Left: &expr.Unary{
			Operator: scanner.Token{
				Type: scanner.Minus,
				Lexeme: "-",
			},
			Right: &expr.Literal{
				Value: 123,
			},
		},
		Operator: scanner.Token{
			Type: scanner.Star,
			Lexeme: "*",
		},
		Right: &expr.Grouping{
			Expression: &expr.Literal{
				Value: 45.67,
			},
		},
	}
	p := AstPrinter{}

	fmt.Println(p.Print(&expr).(string))
}