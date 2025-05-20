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
	err := ReportErrorOnToken(token, format, args...)
	t.errs = append(t.errs, err)
	return err
}

func (t *TokenErrorReporter) HasErrors() bool {
	return len(t.errs) > 0
}

func (t *TokenErrorReporter) Errors() []error {
	return t.errs
}

func (t *TokenErrorReporter) Clear() {
	t.errs = make([]error, 0)
}

func ReportErrorOnToken(token scanner.Token, format string, args ...any) error {
	end := fmt.Sprintf(format, args...)
	err := fmt.Errorf("%v: %s", token.Line, end)
	return err
}
