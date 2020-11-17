package shell

import (
	"context"
	"errors"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSystemShell_Exec(t *testing.T) {
	sh := NewShell(ioutil.Discard, "/")

	tests := []struct {
		name     string
		commands []string
		env      map[string]string
		hasError bool
	}{
		{
			name:     "should return 1,2",
			commands: []string{"echo 1", "echo 2"},
			env:      map[string]string{"demo": "demo"},
		},
		{
			name:     "should fail",
			commands: []string{"exit 10"},
			env:      map[string]string{"demo": "demo"},
			hasError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sh.Exec(context.TODO(), sh.defaultShell, ShScriptConfig{
				Name:     "test",
				Commands: tt.commands,
				Env:      tt.env,
				Output:   true,
			})

			if tt.hasError {
				assert.Error(t, result.Error)
			} else {
				assert.NoError(t, result.Error)
			}
		})
	}
}

func TestExecResult_IsSuccessful(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{name: "Correct execution", err: nil, expected: true},
		{name: "Exit code zero with error", err: errors.New("test"), expected: false},
		{name: "Exit code not zero with error", err: errors.New("test"), expected: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execResult := ExecResult{
				Error: tt.err,
				Name:  "test",
				Time:  time.Second,
			}

			actual := execResult.IsSuccessful()

			assert.Equal(t, tt.expected, actual)
		})
	}
}
