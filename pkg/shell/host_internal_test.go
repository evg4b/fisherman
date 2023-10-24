package shell

import (
	"bytes"
	"context"
	"io"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHost(t *testing.T) {
	var buff bytes.Buffer

	tests := []struct {
		name           string
		shellStr       Strategy
		options        []HostOption
		expectedStdout io.Writer
		expectedStderr io.Writer
		expectedDir    string
		expectedPath   string
		expectedArgs   []string
		expectedEnv    []string
	}{
		{
			name:         "bash without options",
			shellStr:     Bash(),
			options:      []func(str Strategy, host *Host){},
			expectedPath: lookupPath("bash"),
			expectedArgs: []string{"bash"},
		},
		{
			name:     "bash with raw options",
			shellStr: Bash(),
			options: []func(str Strategy, host *Host){
				WithCwd("/demo/test"),
				WithRawArgs([]string{"arg1", "arg2"}),
				WithRawEnv([]string{"VAR1=1", "VAR2=2"}),
			},
			expectedPath: lookupPath("bash"),
			expectedArgs: []string{"arg1", "arg2"},
			expectedEnv:  []string{"VAR1=1", "VAR2=2"},
			expectedDir:  "/demo/test",
		},
		{
			name:     "bash with options",
			shellStr: Bash(),
			options: []func(str Strategy, host *Host){
				WithCwd("/demo/test"),
				WithArgs([]string{"arg1", "arg2"}),
				WithEnv([]string{"VAR1=1", "VAR2=2"}),
			},
			expectedPath: lookupPath("bash"),
			expectedArgs: []string{"arg1", "arg2"},
			expectedEnv:  []string{"VAR1=1", "VAR2=2"},
			expectedDir:  "/demo/test",
		},
		{
			name:     "bash with Stdout and Stderr options",
			shellStr: Bash(),
			options: []func(str Strategy, host *Host){
				WithStdout(&buff),
				WithStderr(&buff),
			},
			expectedPath:   lookupPath("bash"),
			expectedArgs:   []string{"bash"},
			expectedStdout: &buff,
			expectedStderr: &buff,
		},
		{
			name:         "cmd without options",
			shellStr:     Cmd(),
			options:      []func(str Strategy, host *Host){},
			expectedPath: lookupPath("cmd"),
			expectedArgs: []string{"cmd", "/Q", "/D", "/K"},
		},
		{
			name:     "cmd with raw options",
			shellStr: Cmd(),
			options: []func(str Strategy, host *Host){
				WithCwd("/demo/test"),
				WithRawArgs([]string{"arg1", "arg2"}),
				WithRawEnv([]string{"VAR1=1", "VAR2=2"}),
			},
			expectedPath: lookupPath("cmd"),
			expectedArgs: []string{"arg1", "arg2"},
			expectedEnv:  []string{"VAR1=1", "VAR2=2"},
			expectedDir:  "/demo/test",
		},
		{
			name:     "cmd with options",
			shellStr: Cmd(),
			options: []func(str Strategy, host *Host){
				WithCwd("/demo/test"),
				WithArgs([]string{"arg1", "arg2"}),
				WithEnv([]string{"VAR1=1", "VAR2=2"}),
			},
			expectedPath: lookupPath("cmd"),
			expectedArgs: []string{"/Q", "/D", "/K", "arg1", "arg2"},
			expectedEnv:  []string{"VAR1=1", "VAR2=2"},
			expectedDir:  "/demo/test",
		},
		{
			name:     "cmd with Stdout and Stderr options",
			shellStr: Cmd(),
			options: []func(str Strategy, host *Host){
				WithStdout(&buff),
				WithStderr(&buff),
			},
			expectedPath:   lookupPath("cmd"),
			expectedArgs:   []string{"cmd", "/Q", "/D", "/K"},
			expectedStdout: &buff,
			expectedStderr: &buff,
		},
		{
			name:         "powershell without options",
			shellStr:     PowerShell(),
			options:      []func(str Strategy, host *Host){},
			expectedPath: lookupPath("powershell"),
			expectedArgs: []string{"powershell", "-NoProfile", "-NonInteractive", "-NoLogo"},
		},
		{
			name:     "powershell with raw options",
			shellStr: PowerShell(),
			options: []func(str Strategy, host *Host){
				WithCwd("/demo/test"),
				WithRawArgs([]string{"arg1", "arg2"}),
				WithRawEnv([]string{"VAR1=1", "VAR2=2"}),
			},
			expectedPath: lookupPath("powershell"),
			expectedArgs: []string{"arg1", "arg2"},
			expectedEnv:  []string{"VAR1=1", "VAR2=2"},
			expectedDir:  "/demo/test",
		},
		{
			name:     "powershell with options",
			shellStr: PowerShell(),
			options: []func(str Strategy, host *Host){
				WithCwd("/demo/test"),
				WithArgs([]string{"arg1", "arg2"}),
				WithEnv([]string{"VAR1=1", "VAR2=2"}),
			},
			expectedPath: lookupPath("powershell"),
			expectedArgs: []string{"-NoProfile", "-NonInteractive", "-NoLogo", "arg1", "arg2"},
			expectedEnv:  []string{"VAR1=1", "VAR2=2"},
			expectedDir:  "/demo/test",
		},
		{
			name:     "powershell with Stdout and Stderr options",
			shellStr: PowerShell(),
			options: []func(str Strategy, host *Host){
				WithStdout(&buff),
				WithStderr(&buff),
			},
			expectedPath:   lookupPath("powershell"),
			expectedArgs:   []string{"powershell", "-NoProfile", "-NonInteractive", "-NoLogo"},
			expectedStdout: &buff,
			expectedStderr: &buff,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NewHost(context.Background(), tt.shellStr, tt.options...)

			assert.Equal(t, tt.expectedPath, actual.command.Path)
			assert.Equal(t, tt.expectedArgs, actual.command.Args)
			assert.Equal(t, tt.expectedEnv, actual.command.Env)
			assert.Equal(t, tt.expectedDir, actual.command.Dir)
			assert.Equal(t, tt.expectedStderr, actual.command.Stderr)
			assert.Equal(t, tt.expectedStdout, actual.command.Stdout)
		})
	}
}

func lookupPath(name string) string {
	lp, err := exec.LookPath(name)
	if err != nil {
		return name
	}

	return lp
}
