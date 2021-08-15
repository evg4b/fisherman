package shell_test

import (
	"context"
	"fisherman/pkg/shell"
	"fisherman/testing/testutils"
	"io/ioutil"
	"testing"
)

func TestSystemShell_Exec(t *testing.T) {
	sh := shell.NewShell().
		WithWorkingDirectory("/")

	tests := []struct {
		name          string
		commands      []string
		env           map[string]string
		shell         string
		expectedError string
	}{
		{
			name:     "should return 1,2",
			commands: []string{"echo 1", "echo 2"},
			shell:    shell.PlatformDefaultShell,
			env:      map[string]string{"demo": "demo"},
		},
		{
			name:          "should fail wia exit code",
			commands:      []string{"exit 10"},
			env:           map[string]string{"demo": "demo"},
			shell:         shell.PlatformDefaultShell,
			expectedError: "script completed with an error: exit status 10",
		},
		{
			name:          "should fail wia exit code",
			commands:      []string{"echo 1", "echo 2"},
			env:           map[string]string{"demo": "demo"},
			shell:         "unknown-shell",
			expectedError: "failed to get shell configuration: shell 'unknown-shell' is not supported",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sh.Exec(context.TODO(), ioutil.Discard, tt.shell, shell.NewScript(tt.commands).
				SetEnvironmentVariables(tt.env).
				SetDirectory("/"))

			testutils.CheckError(t, tt.expectedError, err)
		})
	}
}

func TestSystemShell_Exec_ShellNotInstalled(t *testing.T) {
	shellName := "test-shell"
	shell.ShellConfigurations = map[string]shell.WrapConfiguration{
		shellName: {Path: shellName},
	}

	sh := shell.NewShell().
		WithWorkingDirectory("/").
		WithDefaultShell(shellName)

	err := sh.Exec(context.TODO(), ioutil.Discard, shellName, shell.NewScript([]string{"echo 1", "echo 2"}).
		SetDirectory("/"))

	testutils.CheckError(t, "failed to get shell configuration: exec: \"test-shell\": executable file not found in %PATH%", err)
}
