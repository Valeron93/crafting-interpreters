package util

import (
	"fmt"

	"github.com/Valeron93/crafting-interpreters/scanner"
)

func ReportErrorOnToken(token scanner.Token, format string, args ...any) error {
	return ReportErrorOnLineAndColumn(token.Line, token.Column, format, args...)
}

func ReportErrorOnLineAndColumn(line int, column int, format string, args ...any) error {
	end := fmt.Sprintf(format, args...)
	return fmt.Errorf("%v:%v: %s", line, column, end)
}
