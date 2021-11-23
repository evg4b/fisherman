package shell

import (
	"context"
	"fisherman/internal/utils"
	pkgutils "fisherman/pkg/utils"
	"fmt"
	"io"
	"os/exec"

	"github.com/go-errors/errors"
)

type SystemShell struct {
	cwd          string
	defaultShell string
	env          []string
}

func NewShell(options ...shellOption) *SystemShell {
	sh := SystemShell{
		defaultShell: PlatformDefaultShell,
		cwd:          "",
		env:          []string{},
	}

	for _, option := range options {
		option(&sh)
	}

	return &sh
}

func (sh *SystemShell) Exec(ctx context.Context, output io.Writer, shell string, script *Script) error {
	config, err := getShellWrapConfiguration(utils.FirstNotEmpty(shell, sh.defaultShell))
	if err != nil {
		return errors.Errorf("failed to get shell configuration: %w", err)
	}

	command := exec.CommandContext(ctx, config.Path, config.Args...) // nolint gosec
	command.Env = pkgutils.MergeEnv(sh.env, script.env)
	command.Dir = utils.FirstNotEmpty(script.dir, sh.cwd)
	// TODO: Add custom encoding for different shell
	command.Stdout = output
	command.Stderr = output

	stdin, err := command.StdinPipe()
	if err != nil {
		return errors.Errorf("failed to get stdin pipe for communication with the process: %w", err)
	}

	go startWriter(stdin, script.commands, config)

	script.duration, err = execWithTime(command.Run)
	if err != nil {
		return errors.Errorf("script completed with an error: %w", err)
	}

	return nil
}

func writeCommandLine(stdin io.Writer, command string) {
	if !utils.IsEmpty(command) {
		fmt.Fprintln(stdin, command)
	}
}

func startWriter(stdin io.WriteCloser, commands []string, config WrapConfiguration) {
	defer stdin.Close()

	writeCommandLine(stdin, config.Init)
	for _, commandLine := range commands {
		writeCommandLine(stdin, commandLine)
		writeCommandLine(stdin, config.PostCommand)
	}

	writeCommandLine(stdin, config.Dispose)
}
