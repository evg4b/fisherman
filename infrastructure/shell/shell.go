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

type ScriptConfig struct {
	Name     string
	Commands []string          `yaml:"commands,omitempty"`
	Env      map[string]string `yaml:"env,omitempty"`
	Output   bool              `yaml:"output,omitempty"`
}

type SystemShell struct {
	output io.Writer
}

func NewShell(output io.Writer) *SystemShell {
	return &SystemShell{
		output: output,
	}
}

func (sh *SystemShell) Exec(ctx context.Context, script ScriptConfig) ExecResult {
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
