package shell

import (
	"context"
	"fisherman/internal/utils"
	"fmt"
	"io"
	"os"
)

type SystemShell struct {
	cwd          string
	defaultShell string
}

func NewShell(output io.Writer, cwd, defaultShell string) *SystemShell {
	return &SystemShell{cwd, utils.GetOrDefault(defaultShell, DefaultShell)}
}

func (sh *SystemShell) Exec(ctx context.Context, output io.Writer, shell string, script *Script) error {
	envList := os.Environ()
	for key, value := range script.env {
		envList = append(envList, fmt.Sprintf("%s=%s", key, value))
	}

	command, err := CommandFactory(ctx, utils.GetOrDefault(shell, sh.defaultShell))
	if err != nil {
		return err
	}

	command.Env = envList
	command.Dir = utils.GetOrDefault(script.dir, sh.cwd)
	command.Stdout = output
	command.Stderr = output

	stdin, err := command.StdinPipe()
	if err != nil {
		panic(err)
	}

	go func() {
		defer stdin.Close()
		for _, commandLine := range script.commands {
			fmt.Fprintln(stdin, commandLine)
		}
	}()

	return command.Run()
}
