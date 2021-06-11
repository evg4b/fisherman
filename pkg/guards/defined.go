package guards

import "github.com/go-errors/errors"

func ShouldBeDefined(object interface{}, message string) {
	if object == nil {
		panic(errors.New(message))
	}
}
