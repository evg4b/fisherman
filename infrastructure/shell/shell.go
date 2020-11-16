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
	output       io.Writer
	cwd          string
	defaultShell string
}

func NewShell(output io.Writer, cwd string) *SystemShell {
	return &SystemShell{output, cwd, DefaultShell}
}

func (sh *SystemShell) Exec(ctx context.Context, shell string, script ShScriptConfig) ExecResult {
	envList := os.Environ()
	for key, value := range script.Env {
		envList = append(envList, fmt.Sprintf("%s=%s", key, value))
	}

	command, err := CommandFactory(ctx, utils.GetOrDefault(shell, sh.defaultShell))
	if err != nil {
		return ExecResult{
			Error: err,
			Time:  time.Duration(0),
		}
	}

	command.Env = envList
	command.Dir = utils.GetOrDefault(script.Dir, sh.cwd)

	if script.Output {
		prefix := fmt.Sprintf("%s |", script.Name)
		output := prefixer.New(sh.output, func() string {
			return prefix
		})
		command.Stdout = output
		command.Stderr = output
	}

	stdin, err := command.StdinPipe()
	if err != nil {
		panic(err)
	}

	go func() {
		defer stdin.Close()
		for _, commandLine := range script.Commands {
			fmt.Fprintln(stdin, commandLine)
		}
	}()

	duration, err := utils.ExecWithTime(command.Run)

	return ExecResult{
		Error: err,
		Time:  duration,
		Name:  script.Name,
	}
}
