package util

import (
	"fmt"

	"github.com/Valeron93/crafting-interpreters/scanner"
)

type TokenErrorReporter struct {
	errs []error
}

func NewTokenErrorReporter() TokenErrorReporter {
	return TokenErrorReporter{
		errs: make([]error, 0),
	}
}

func (t *TokenErrorReporter) Report(token scanner.Token, format string, args ...any) error {
	end := fmt.Sprintf(format, args...)
	return fmt.Errorf("line %v: %s", token.Line, end)
}

func (t *TokenErrorReporter) HasErrors() bool {
	return len(t.errs) > 0
}

func (t *TokenErrorReporter) Errors() []error {
	return t.errs
}
