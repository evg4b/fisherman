package validation

import "github.com/hashicorp/go-multierror"

func IsValidationError(err error) bool {
	if multierror, ok := err.(*multierror.Error); ok {
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
	_, ok := err.(*Error)

	return ok
}
