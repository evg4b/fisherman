package shell

import (
	"context"
	"os/exec"
)

type BashStrategy struct{}

func Bash() *BashStrategy {
	return &BashStrategy{}
}

func (s *BashStrategy) GetName() string {
	return "bash"
}

func (s *BashStrategy) GetCommand(ctx context.Context) *exec.Cmd {
	args := s.ArgsWrapper([]string{})

	return exec.CommandContext(ctx, "bash", args...)
}

func (s *BashStrategy) ArgsWrapper(args []string) []string {
	defaultArgs := []string{"-i"}

	return append(defaultArgs, args...)
}

func (s *BashStrategy) EnvWrapper(env []string) []string {
	return env
}
