package shell

import (
	"bytes"
	"context"
	"fisherman/utils"
	"fmt"
	"os"
	"time"
)

type ScriptConfig struct {
	Name     string
	Commands []string          `yaml:"commands,omitempty"`
	Env      map[string]string `yaml:"env,omitempty"`
	Output   bool              `yaml:"output,omitempty"`
}

type SystemShell struct {
}

func NewShell() *SystemShell {
	return &SystemShell{}
}

func (*SystemShell) Exec(ctx context.Context, script ScriptConfig) ExecResult {
	var stdout bytes.Buffer

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
		command.Stdout = &stdout
		command.Stderr = &stdout
	}

	duration, err := utils.ExecWithTime(command.Run)

	return ExecResult{
		Output:   stdout.String(),
		Error:    err,
		ExitCode: command.ProcessState.ExitCode(),
		Time:     duration,
		Name:     script.Name,
	}
}
