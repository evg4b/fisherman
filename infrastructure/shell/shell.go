package shell

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type shellCommandFactory = func(commands []string) (*exec.Cmd, error)

var CommandFactory shellCommandFactory

type SystemShell struct {
}

func NewShell() *SystemShell {
	return &SystemShell{}
}

func (*SystemShell) Exec(commands []string, env *map[string]string, paths []string) (string, string, int, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if len(paths) > 0 {
		(*env)["PATH"] = makePathVariable(paths)
	}

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

	return stdout.String(), stderr.String(), cmd.ProcessState.ExitCode(), err
}

func makePathVariable(paths []string) string {
	pathsList := make([]string, 10)
	path, exists := os.LookupEnv("PATH")

	if exists {
		pathsList = append(pathsList, strings.Split(path, PathVariableSeparator)...)
	}

	if len(paths) > 0 {
		pathsList = append(pathsList, paths...)
	}

	return strings.Join(pathsList, PathVariableSeparator)
}
