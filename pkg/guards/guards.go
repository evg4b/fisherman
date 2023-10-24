package guards

import (
	"reflect"
	"strings"

	"github.com/go-errors/errors"
)

// ShouldBeDefined checks that the variable is not nil. In case when a variable is nill, it panic.
func ShouldBeDefined(object any, message string) {
	if object == nil || (reflect.ValueOf(object).Kind() == reflect.Ptr && reflect.ValueOf(object).IsNil()) {
		panic(errors.New(message))
	}
}

// NoError panic when passed error is not nil.
func NoError(err error) {
	if err != nil {
		panic(err)
	}
}

// ShouldBeNotEmpty checks that the string is not empty and should contins not witespace symbols.
func ShouldBeNotEmpty(object string, message string) {
	if len(strings.TrimSpace(object)) == 0 {
		panic(errors.New(message))
	}
}
