package shell_test

import (
	"context"
	"fisherman/pkg/shell"
	"fisherman/testing/testutils"
	"io/ioutil"
	"runtime"
	"testing"
)

func TestSystemShell_Exec(t *testing.T) {
	sh := shell.NewShell(shell.WithWorkingDirectoryOld("/"))

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
			script := shell.NewScript(tt.commands).
				SetEnvironmentVariables(tt.env).
				SetDirectory("/")

			err := sh.Exec(context.TODO(), ioutil.Discard, tt.shell, script)

			testutils.CheckError(t, tt.expectedError, err)
		})
	}
}

func TestSystemShell_Exec_ShellNotInstalled(t *testing.T) {
	shellName := "test-shell"
	shell.ShellConfigurations = map[string]shell.WrapConfiguration{
		shellName: {Path: shellName},
	}

	sh := shell.NewShell(
		shell.WithWorkingDirectoryOld("/"),
		shell.WithDefaultShell(shellName),
	)

	script := shell.NewScript([]string{"echo 1", "echo 2"}).
		SetDirectory("/")

	err := sh.Exec(context.TODO(), ioutil.Discard, shellName, script)

	if runtime.GOOS == "windows" {
		testutils.CheckError(t, "failed to get shell configuration: exec: \"test-shell\": executable file not found in %PATH%", err)
	} else {
		testutils.CheckError(t, "failed to get shell configuration: exec: \"test-shell\": executable file not found in $PATH", err)
	}
}
