package guards_test

import (
	"testing"

	. "github.com/evg4b/fisherman/pkg/guards"

	"github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"
)

func TestShouldBeDefined(t *testing.T) {
	var nilPointer *struct{}

	tests := []struct {
		name    string
		object  any
		message string
		err     string
	}{
		{
			name:    "should not panic for zero",
			message: "unknown",
			object:  0,
		},
		{
			name:    "should not panic for defined empty string",
			message: "unknown",
			object:  "",
		},
		{
			name:    "should not panic for defined empty struct",
			message: "unknown",
			object:  struct{}{},
		},
		{
			name:    "should not panic for defined pointer to struct",
			message: "unknown",
			object:  &struct{}{},
		},
		{
			name:    "should panic for defined pointer to struct",
			message: "value is nil",
			object:  nilPointer,
			err:     "value is nil",
		},
		{
			name:    "should not panic for defined empty slice",
			message: "unknown",
			object:  []string{},
		},
		{
			name:    "should panic for nil value",
			message: "value is nil",
			err:     "value is nil",
			object:  nil,
		},
	}

	for _, testCase := range tests {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			if len(testCase.err) > 0 {
				assert.PanicsWithError(t, testCase.err, func() {
					ShouldBeDefined(testCase.object, testCase.message)
				})
			} else {
				assert.NotPanics(t, func() {
					ShouldBeDefined(testCase.object, testCase.message)
				})
			}
		})
	}
}

func TestNoError(t *testing.T) {
	t.Run("panic for error", func(t *testing.T) {
		err := errors.New("Test err")
		assert.PanicsWithError(t, err.Error(), func() {
			NoError(err)
		})
	})
	t.Run("does not panic for nil", func(t *testing.T) {
		assert.NotPanics(t, func() {
			NoError(nil)
		})
	})
}

func TestShouldBeNotEmpty(t *testing.T) {
	tests := []struct {
		name   string
		object string
		err    string
	}{
		{
			name:   "should not panic for string with numbers",
			object: "0",
		},
		{
			name:   "should panic for empty string",
			object: "",
			err:    "string is empty",
		},
		{
			name:   "should panic for tabs",
			object: "\t\t",
			err:    "string is empty",
		},
		{
			name:   "should panic for carret symbol",
			object: "\r\r",
			err:    "string is empty",
		},
		{
			name:   "should panic for spaces",
			object: "   ",
			err:    "string is empty",
		},
		{
			name:   "should not panic for string with witespace symbols",
			object: "\n\r\t  not empty",
		},
		{
			name:   "should panic for mixed content",
			err:    "string is empty",
			object: " \t \n   ",
		},
	}

	for _, testCase := range tests {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			if len(testCase.err) > 0 {
				assert.PanicsWithError(t, testCase.err, func() {
					ShouldBeNotEmpty(testCase.object, "string is empty")
				})
			} else {
				assert.NotPanics(t, func() {
					ShouldBeNotEmpty(testCase.object, "string is empty")
				})
			}
		})
	}
}
