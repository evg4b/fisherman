//go:build debug
// +build debug

package log

import (
	"github.com/go-errors/errors"
	"github.com/hashicorp/go-multierror"
)

func DumpError(err error) {
	if withStack, ok := err.(*errors.Error); ok {
		printWithStackError(withStack)
	}

	if multiError, ok := err.(*multierror.Error); ok {
		for _, err := range multiError.Errors {
			DumpError(err)
		}
	}
}

func printWithStackError(err *errors.Error) {
	Errorf("===> [debug]: %s", err.ErrorStack())
}
