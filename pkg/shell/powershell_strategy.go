package shell

import (
	"context"
	"os/exec"
)

type PowershellStrategy struct{}

func PowerShell() *PowershellStrategy {
	return &PowershellStrategy{}
}

func (c *PowershellStrategy) GetCommand(ctx context.Context) *exec.Cmd {
	args := c.ArgsWrapper([]string{})

	return exec.CommandContext(ctx, "powershell", args...)
}

func (c *PowershellStrategy) ArgsWrapper(args []string) []string {
	defaultArgs := []string{"-NoProfile", "-NonInteractive", "-NoLogo"}

	return append(defaultArgs, args...)
}

func (c *PowershellStrategy) EnvWrapper(env []string) []string {
	return env
}
