package shell

import (
	"context"
	"os/exec"
)

type CmdStrategy struct{}

const cmdBin = "cmd"

func Cmd() *CmdStrategy {
	return &CmdStrategy{}
}

func (s *CmdStrategy) GetName() string {
	return cmdBin
}

func (s *CmdStrategy) GetCommand(ctx context.Context) *exec.Cmd {
	args := s.ArgsWrapper([]string{})

	return exec.CommandContext(ctx, cmdBin, args...)
}

func (s *CmdStrategy) ArgsWrapper(args []string) []string {
	return append([]string{"/Q", "/D", "/K"}, args...)
}

func (s *CmdStrategy) EnvWrapper(env []string) []string {
	return env
}
