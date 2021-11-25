package shell

import (
	"context"
	"os/exec"
)

type CmdStrategy struct{}

func Cmd() *CmdStrategy {
	return &CmdStrategy{}
}

func (c *CmdStrategy) GetCommand(ctx context.Context) *exec.Cmd {
	args := c.ArgsWrapper([]string{})

	return exec.CommandContext(ctx, "cmd", args...)
}

func (c *CmdStrategy) ArgsWrapper(args []string) []string {
	defaultArgs := []string{"/Q", "/D", "/K"}

	return append(defaultArgs, args...)
}

func (c *CmdStrategy) EnvWrapper(env []string) []string {
	return env
}
