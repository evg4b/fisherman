package validation

import (
	"errors"

	"github.com/hashicorp/go-multierror"
)

func IsValidationError(err error) bool {
	var multierror *multierror.Error
	if errors.As(err, &multierror) {
		for _, e := range multierror.Errors {
			if !isValidationErrorInternal(e) {
				return false
			}
		}

		return true
	}

	return isValidationErrorInternal(err)
}

func isValidationErrorInternal(err error) bool {
	var validationError *Error

	return errors.As(err, &validationError)
}
