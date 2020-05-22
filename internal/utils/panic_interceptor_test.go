package utils_test

import (
	"errors"
	"testing"

	. "github.com/evg4b/fisherman/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestPanicInterceptor(t *testing.T) {
	tests := []struct {
		name           string
		panicData      any
		shouldBeCalled bool
	}{
		{
			name:           "intercepts panic and return with exit code 3",
			panicData:      "test panic",
			shouldBeCalled: true,
		},
		{
			name:           "intercepts panic and return with exit code 0",
			panicData:      errors.New("test error"),
			shouldBeCalled: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			called := false

			assert.NotPanics(t, func() {
				defer PanicInterceptor(func(data any) {
					called = true
					assert.Equal(t, tt.panicData, data)
				})

				panic(tt.panicData)
			})

			assert.Equal(t, tt.shouldBeCalled, called)
		})
	}
}
