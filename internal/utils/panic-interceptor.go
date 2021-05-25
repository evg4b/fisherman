package utils

import (
	"fisherman/infrastructure/log"
)

func PanicInterceptor(exit func(code int), exitCode int) {
	if err := recover(); err != nil {
		log.Errorf("Fatal error: %s", err)
		exit(exitCode)
	}
}
