package validation

import (
	"errors"

	"github.com/hashicorp/go-multierror"
)

func IsValidationError(err error) bool {
	var multiError *multierror.Error
	if errors.As(err, &multiError) {
		for _, e := range multiError.Errors {
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
