package shell

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"
)

type ExecResult struct {
	Output   string
	ExitCode int
	Error    error
	Time     time.Duration
}

type SystemShell struct {
}

func NewShell() *SystemShell {
	return &SystemShell{}
}

func (*SystemShell) Exec(commands []string, env *map[string]string, output bool) ExecResult {
	var stdout bytes.Buffer

	envList := os.Environ()
	for key, value := range *env {
		envList = append(envList, fmt.Sprintf("%s=%s", key, value))
	}

	command, err := CommandFactory(context.TODO(), commands)
	if err != nil {
		return ExecResult{
			Error:    err,
			ExitCode: -1,
			Time:     time.Duration(0),
		}
	}

	command.Env = envList
	if output {
		command.Stdout = &stdout
		command.Stderr = &stdout
	}

	start := time.Now()
	err = command.Run()
	executionTime := time.Since(start)

	return ExecResult{
		Output:   stdout.String(),
		Error:    err,
		ExitCode: command.ProcessState.ExitCode(),
		Time:     executionTime,
	}
}
