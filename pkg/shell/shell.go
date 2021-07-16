package shell

import (
	"context"
	"fisherman/internal/utils"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type SystemShell struct {
	cwd          string
	defaultShell string
}

func NewShell() *SystemShell {
	return &SystemShell{defaultShell: PlatformDefaultShell}
}

func (sh *SystemShell) WithWorkingDirectory(cwd string) *SystemShell {
	sh.cwd = cwd

	return sh
}

func (sh *SystemShell) WithDefaultShell(defaultShell string) *SystemShell {
	sh.defaultShell = defaultShell

	return sh
}

func (sh *SystemShell) Exec(ctx context.Context, output io.Writer, shell string, script *Script) error {
	envList := os.Environ()
	for key, value := range script.env {
		envList = append(envList, fmt.Sprintf("%s=%s", key, value))
	}

	config, err := getShellWrapConfiguration(utils.GetOrDefault(shell, sh.defaultShell))
	if err != nil {
		return err
	}

	command := exec.CommandContext(ctx, config.Path, config.Args...) // nolint gosec
	command.Env = envList
	command.Dir = utils.GetOrDefault(script.dir, sh.cwd)
	command.Stdout = output
	command.Stderr = output

	stdin, err := command.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		defer stdin.Close()
		writeCommandLine(stdin, config.Init)
		for _, commandLine := range script.commands {
			writeCommandLine(stdin, commandLine)
			writeCommandLine(stdin, config.PostCommand)
		}
		writeCommandLine(stdin, config.Dispose)
	}()

	script.duration, err = utils.ExecWithTime(command.Run)
	if err != nil {
		return err
	}

	return err
}

func writeCommandLine(stdin io.Writer, command string) {
	if !utils.IsEmpty(command) {
		fmt.Fprintln(stdin, command)
	}
}
