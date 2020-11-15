package shell

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestExecResult_IsCanceled(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{name: "No error", err: nil, expected: false},
		{name: "Canceled error", err: context.Canceled, expected: true},
		{name: "Other error", err: errors.New("test"), expected: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExecResult{Error: tt.err}

			isCanceled := result.IsCanceled()

			assert.Equal(t, tt.expected, isCanceled)
		})
	}
}
