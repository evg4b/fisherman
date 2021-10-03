package validation

import "fmt"

// Error is base error.
type Error struct {
	prefix  string
	message string
}

func (err *Error) Error() string {
	if len(err.prefix) > 0 {
		return fmt.Sprintf("[%s] %s", err.prefix, err.message)
	}

	return err.message
}

func Errorf(prefix, message string, a ...interface{}) error {
	return &Error{
		prefix:  prefix,
		message: fmt.Sprintf(message, a...),
	}
}
