package shell

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type ShellCommandFactory = func(commands []string) (*exec.Cmd, error)

var CommandFactory ShellCommandFactory

type SystemShell struct {
}

func NewShell() *SystemShell {
	return &SystemShell{}
}

func (*SystemShell) Exec(commands []string, env *map[string]string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	envList := os.Environ()
	for key, value := range *env {
		envList = append(envList, fmt.Sprintf("%s=%s", key, value))
	}

	cmd, err := CommandFactory(commands)
	if err == nil {
		cmd.Env = envList
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err = cmd.Run()
	}

	return stdout.String(), stderr.String(), err
}
