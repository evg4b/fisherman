package shell

import (
	"bytes"
	"fmt"
	"os"
)

type SystemShell struct {
}

func NewShell() *SystemShell {
	return &SystemShell{}
}

func (*SystemShell) Exec(commands []string, env *map[string]string, output bool) (string, int, error) {
	var stdout bytes.Buffer

	envList := os.Environ()
	for key, value := range *env {
		envList = append(envList, fmt.Sprintf("%s=%s", key, value))
	}

	command, err := CommandFactory(commands)
	if err != nil {
		return "", -1, err
	}

	command.Env = envList
	if output {
		command.Stdout = &stdout
		command.Stderr = &stdout
	}

	err = command.Run()

	return stdout.String(), command.ProcessState.ExitCode(), err
}
