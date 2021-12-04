package shell

import (
	"context"
	"os/exec"
)

type BashStrategy struct{}

const bashBin = "bash"

func Bash() *BashStrategy {
	return &BashStrategy{}
}

func (s *BashStrategy) GetName() string {
	return bashBin
}

func (s *BashStrategy) GetCommand(ctx context.Context) *exec.Cmd {
	args := s.ArgsWrapper([]string{})

	return exec.CommandContext(ctx, bashBin, args...)
}

func (s *BashStrategy) ArgsWrapper(args []string) []string {
	return append([]string{"-i"}, args...)
}

func (s *BashStrategy) EnvWrapper(env []string) []string {
	return env
}
