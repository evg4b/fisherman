package shell_test

import (
	"context"
	"fisherman/pkg/shell"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystemShell_Exec(t *testing.T) {
	sh := shell.NewShell().
		WithWorkingDirectory("/")

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
			result := sh.Exec(context.TODO(), ioutil.Discard, shell.DefaultShell, shell.NewScript(tt.commands).
				SetEnvironmentVariables(tt.env).
				SetDirectory("/"))

			if tt.hasError {
				assert.Error(t, result)
			} else {
				assert.NoError(t, result)
			}
		})
	}
}
