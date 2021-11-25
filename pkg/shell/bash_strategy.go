package shell

import (
	"context"
	"os/exec"
)

type BashStrategy struct{}

func Bash() *BashStrategy {
	return &BashStrategy{}
}

func (c *BashStrategy) GetCommand(ctx context.Context) *exec.Cmd {
	args := c.ArgsWrapper([]string{})

	return exec.CommandContext(ctx, "bash", args...)
}

func (c *BashStrategy) ArgsWrapper(args []string) []string {
	defaultArgs := []string{"-i"}

	return append(defaultArgs, args...)
}

func (c *BashStrategy) EnvWrapper(env []string) []string {
	return env
}
