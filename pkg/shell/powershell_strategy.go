package shell

import (
	"context"
	"os/exec"
)

const powershellBin = "powershell"

type PowershellStrategy struct{}

func PowerShell() *PowershellStrategy {
	return &PowershellStrategy{}
}

func (s *PowershellStrategy) GetName() string {
	return powershellBin
}

func (s *PowershellStrategy) GetCommand(ctx context.Context) *exec.Cmd {
	args := s.ArgsWrapper([]string{})

	return exec.CommandContext(ctx, powershellBin, args...)
}

func (s *PowershellStrategy) ArgsWrapper(args []string) []string {
	return append([]string{"-NoProfile", "-NonInteractive", "-NoLogo"}, args...)
}

func (s *PowershellStrategy) EnvWrapper(env []string) []string {
	return env
}
