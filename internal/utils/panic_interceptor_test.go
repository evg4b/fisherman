package utils_test

import (
	"bytes"
	"errors"
	. "fisherman/internal/utils"
	"fisherman/pkg/log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPanicInterceptor(t *testing.T) {
	t.Run("corrected exit code", func(t *testing.T) {
		tests := []struct {
			name     string
			exitCode int
		}{
			{
				name:     "intercepts panic and return with exit code 3",
				exitCode: 3,
			},
			{
				name:     "intercepts panic and return with exit code 0",
				exitCode: 0,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				called := false

				assert.NotPanics(t, func() {
					defer PanicInterceptor(func(code int) {
						assert.Equal(t, tt.exitCode, code)
						called = true
					}, tt.exitCode)
					panic("test panic")
				})

				assert.True(t, called)
			})
		}
	})

	t.Run("error dump", func(t *testing.T) {
		buffer := bytes.NewBufferString("")

		log.SetOutput(buffer)

		called := false
		assert.NotPanics(t, func() {
			defer PanicInterceptor(func(int) { called = true }, 3)
			panic(errors.New("test panic"))
		})

		assert.True(t, called)
		assert.Equal(t, buffer.String(), "Fatal error: test panic\n")
	})
}
