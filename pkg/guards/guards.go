package guards

import "github.com/go-errors/errors"

// ShouldBeDefined checks that the variable is not nil. In case when a variable is nill, it panic.
func ShouldBeDefined(object interface{}, message string) {
	if object == nil {
		panic(errors.New(message))
	}
}

// NoError panic when passed error is not nil.
func NoError(err error) {
	if err != nil {
		panic(err)
	}
}
