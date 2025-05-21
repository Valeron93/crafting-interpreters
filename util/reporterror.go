package util

import (
	"fmt"

	"github.com/Valeron93/crafting-interpreters/scanner"
)

func ReportErrorOnToken(token scanner.Token, format string, args ...any) error {
	end := fmt.Sprintf(format, args...)
	err := fmt.Errorf("%v:%v: %s", token.Line, token.Column, end)
	return err
}
