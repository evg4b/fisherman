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
		{name: "1", err: nil, expected: false},
		{name: "2", err: context.Canceled, expected: true},
		{name: "3", err: errors.New("test"), expected: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExecResult{Error: tt.err}

			isCanceled := result.IsCanceled()

			assert.Equal(t, tt.expected, isCanceled)
		})
	}
}
