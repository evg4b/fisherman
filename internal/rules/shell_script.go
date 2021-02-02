package rules

import (
	"fisherman/infrastructure/shell"
	"fisherman/internal"
	"fisherman/internal/prefixwriter"
	"fmt"
	"io"
	"io/ioutil"
)

const ShellScriptType = "shell-script"

type ShellScript struct {
	BaseRule `mapstructure:",squash"`
	Name     string            `mapstructure:"name"`
	Shell    string            `mapstructure:"shell"`
	Commands []string          `mapstructure:"commands"`
	Env      map[string]string `mapstructure:"env"`
	Output   bool              `mapstructure:"output"`
	Dir      string            `mapstructure:"dir"`
}

func (config *ShellScript) GetPosition() byte {
	return Scripts
}

func (config *ShellScript) Check(ctx internal.ExecutionContext, output io.Writer) error {
	return ctx.Shell().
		Exec(ctx, formatOutput(output, config), config.Shell, shell.ShScript{
			Commands: config.Commands,
			Env:      config.Env,
			Dir:      config.Dir,
		})
}

func formatOutput(output io.Writer, config *ShellScript) io.Writer {
	if config.Output {
		return prefixwriter.New(output, fmt.Sprintf("[%s] ", config.Name))
	}

	return ioutil.Discard
}
