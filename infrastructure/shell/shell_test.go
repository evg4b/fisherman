package shell

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystemShell_Exec(t *testing.T) {
	sh := NewShell()

	tests := []struct {
		name           string
		commands       []string
		env            *map[string]string
		expectedStdout string
		err            error
	}{
		{
			name:           "should return 1,2",
			commands:       []string{"echo 1", "echo 2"},
			env:            &map[string]string{"demo": "demo"},
			err:            nil,
			expectedStdout: fmt.Sprintf("1%s2%s", LineBreak, LineBreak),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, exitCode, err := sh.Exec(tt.commands, tt.env, []string{})
			assert.Equal(t, tt.expectedStdout, stdout)
			assert.Equal(t, 0, exitCode)
			assert.NoError(t, err)
		})
	}
}

func Test_makePathVariable(t *testing.T) {
	tests := []struct {
		name         string
		paths        []string
		expectedPath string
		sysPath      string
	}{
		{
			name:         "test",
			sysPath:      "",
			paths:        []string{},
			expectedPath: "",
		},
		{
			name:         "test",
			sysPath:      "/root/test",
			paths:        []string{},
			expectedPath: "/root/test",
		},
		{
			name:         "test",
			sysPath:      "/root/test",
			paths:        []string{"/bin", "/usr/root"},
			expectedPath: fmt.Sprintf("/root/test%s/bin%s/usr/root", PathVariableSeparator, PathVariableSeparator),
		},
		{
			name:         "test",
			sysPath:      fmt.Sprintf("/bin%s/usr/root", PathVariableSeparator),
			paths:        []string{"/root/test"},
			expectedPath: fmt.Sprintf("/bin%s/usr/root%s/root/test", PathVariableSeparator, PathVariableSeparator),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("PATH", tt.sysPath)
			path := makePathVariable(tt.paths)
			assert.Equal(t, tt.expectedPath, path)
		})
	}
}
