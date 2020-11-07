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

func (*SystemShell) Exec(commands []string, env *map[string]string) (string, int, error) {
	var stdout bytes.Buffer

	envList := os.Environ()
	for key, value := range *env {
		envList = append(envList, fmt.Sprintf("%s=%s", key, value))
	}

	cmd, err := CommandFactory(commands)
	if err == nil {
		cmd.Env = envList
		cmd.Stdout = &stdout
		cmd.Stderr = &stdout
		err = cmd.Run()
	}

	return stdout.String(), cmd.ProcessState.ExitCode(), err
}
