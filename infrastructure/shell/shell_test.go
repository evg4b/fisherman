package shell

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSystemShell_Exec(t *testing.T) {
	sh := NewShell()

	tests := []struct {
		name           string
		commands       []string
		env            map[string]string
		expectedStdout string
		exitCode       int
	}{
		{
			name:           "should return 1,2",
			commands:       []string{"echo 1", "echo 2"},
			env:            map[string]string{"demo": "demo"},
			expectedStdout: fmt.Sprintf("1%s2%s", LineBreak, LineBreak),
			exitCode:       0,
		},
		{
			name:     "should return 1,2",
			commands: []string{"demo"},
			env:      map[string]string{"demo": "demo"},
			exitCode: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sh.Exec(context.TODO(), ScriptConfig{
				Name:     "test",
				Commands: tt.commands,
				Env:      tt.env,
				Output:   true,
			})

			assert.Equal(t, tt.exitCode, result.ExitCode)
			if tt.exitCode == 0 {
				assert.Equal(t, tt.expectedStdout, result.Output)
				assert.NoError(t, result.Error)
			} else {
				assert.Error(t, result.Error)
			}
		})
	}
}

func TestExecResult_IsSuccessful(t *testing.T) {
	tests := []struct {
		name     string
		exitCode int
		err      error
		expected bool
	}{
		{name: "Correct execution", exitCode: 0, err: nil, expected: true},
		{name: "Exit code -1 zero", exitCode: -1, err: nil, expected: false},
		{name: "Exit code 1 zero", exitCode: 1, err: nil, expected: false},
		{name: "Exit code zero with error", exitCode: 0, err: errors.New("test"), expected: false},
		{name: "Exit code not zero with error", exitCode: 1, err: errors.New("test"), expected: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execResult := ExecResult{
				ExitCode: tt.exitCode,
				Error:    tt.err,
				Name:     "test",
				Output:   "",
				Time:     time.Second,
			}

			actual := execResult.IsSuccessful()

			assert.Equal(t, tt.expected, actual)
		})
	}
}
