package utils

import (
	"fisherman/pkg/log"
)

func PanicInterceptor(exit func(code int), exitCode int) {
	if recovered := recover(); recovered != nil {
		log.Errorf("Fatal error: %s", recovered)
		if err, ok := recovered.(error); ok {
			log.DumpError(err)
		}

		exit(exitCode)
	}
}
