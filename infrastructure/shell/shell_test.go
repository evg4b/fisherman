package shell_test

import (
	"fisherman/infrastructure/shell"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystemShell_Exec(t *testing.T) {
	sh := shell.NewShell()

	tests := []struct {
		name           string
		commands       []string
		env            *map[string]string
		expectedStdout string
		expectedStderr string
		err            error
	}{
		{
			name: "should return 1,2",
			commands: []string{
				"echo 1",
				"echo 2",
			},
			env: &map[string]string{
				"demo": "demo",
			},
			err:            nil,
			expectedStderr: "",
			expectedStdout: fmt.Sprintf("1%s2%s", shell.LineBreak, shell.LineBreak),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, exitCode, err := sh.Exec(tt.commands, tt.env, []string{})
			assert.Equal(t, tt.expectedStdout, stdout)
			assert.Equal(t, tt.expectedStderr, stderr)
			assert.Equal(t, 0, exitCode)
			assert.NoError(t, err)
		})
	}
}
