package shell

import (
	"context"
	"fisherman/utils"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/egymgmbh/go-prefix-writer/prefixer"
)

type ShScriptConfig struct {
	Name     string
	Commands []string
	Env      map[string]string
	Output   bool
	Dir      string
}

type SystemShell struct {
	output io.Writer
}

func NewShell(output io.Writer) *SystemShell {
	return &SystemShell{output}
}

func (sh *SystemShell) Exec(ctx context.Context, script ShScriptConfig) ExecResult {
	envList := os.Environ()
	for key, value := range script.Env {
		envList = append(envList, fmt.Sprintf("%s=%s", key, value))
	}

	command, err := CommandFactory(ctx, script.Commands)
	if err != nil {
		return ExecResult{
			Error:    err,
			ExitCode: -1,
			Time:     time.Duration(0),
		}
	}

	command.Env = envList
	if utils.IsNotEmpty(script.Dir) {
		command.Dir = script.Dir
	}

	if script.Output {
		prefix := fmt.Sprintf("%s |", script.Name)
		output := prefixer.New(sh.output, func() string {
			return prefix
		})
		command.Stdout = output
		command.Stderr = output
	}

	duration, err := utils.ExecWithTime(command.Run)

	return ExecResult{
		Error:    err,
		ExitCode: command.ProcessState.ExitCode(),
		Time:     duration,
		Name:     script.Name,
	}
}
