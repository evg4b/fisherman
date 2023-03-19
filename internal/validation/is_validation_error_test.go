package validation_test

import (
	"errors"
	. "fisherman/internal/validation"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
)

func TestIsValidationError(t *testing.T) {
	validationError := Errorf("test", "test")
	notValidationError := errors.New("test")

	var singleNotValidationError *multierror.Error
	singleNotValidationError = multierror.Append(singleNotValidationError, notValidationError)

	var singleValidationError *multierror.Error
	singleValidationError = multierror.Append(singleValidationError, validationError)

	var multiNotValidationErrors *multierror.Error
	multiNotValidationErrors = multierror.Append(multiNotValidationErrors, validationError)
	multiNotValidationErrors = multierror.Append(multiNotValidationErrors, validationError)

	var multiValidationErrors *multierror.Error
	multiValidationErrors = multierror.Append(multiValidationErrors, validationError)
	multiValidationErrors = multierror.Append(multiValidationErrors, notValidationError)

	testCases := []struct {
		name     string
		err      error
		expected bool
	}{
		{name: "Validation error", err: validationError, expected: true},
		{name: "Not validation error", err: notValidationError, expected: false},
		{name: "Multierror with single not validation error", err: singleNotValidationError, expected: false},
		{name: "Multierror with single validation error", err: singleValidationError, expected: true},
		{name: "Multierror with not validation error", err: multiNotValidationErrors, expected: true},
		{name: "Multierror without not validation error", err: multiValidationErrors, expected: false},
	}
	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			actual := IsValidationError(testCase.err)

			assert.Equal(t, testCase.expected, actual)
		})
	}
}
