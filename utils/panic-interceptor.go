package utils

import (
	"fisherman/infrastructure/log"
)

// PanicInterceptor intercept panic and call exit function with exit code
func PanicInterceptor(exit func(code int), exitCode int) {
	if err := recover(); err != nil {
		log.Errorf("Fatal error: %s", err)
		exit(exitCode)
	}
}
