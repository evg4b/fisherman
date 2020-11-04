package shell

import (
	"bytes"
	"fisherman/utils"
	"fmt"
	"os"
	"strings"
)

type SystemShell struct {
}

func NewShell() *SystemShell {
	return &SystemShell{}
}

func (*SystemShell) Exec(commands []string, env *map[string]string, paths []string) (string, int, error) {
	var stdout bytes.Buffer

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
		cmd.Stderr = &stdout
		err = cmd.Run()
	}

	return stdout.String(), cmd.ProcessState.ExitCode(), err
}

func makePathVariable(paths []string) string {
	pathsList := []string{}
	path, exists := os.LookupEnv("PATH")

	if exists {
		pathsList = append(pathsList, strings.Split(path, PathVariableSeparator)...)
	}

	if len(paths) > 0 {
		pathsList = append(pathsList, paths...)
	}

	filtered := utils.Filter(pathsList, utils.IsNotEmpty)

	return strings.Join(filtered, PathVariableSeparator)
}
