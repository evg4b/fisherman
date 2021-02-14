package rules

import (
	"fisherman/infrastructure/shell"
	"fisherman/internal"
	"fisherman/internal/prefixwriter"
	"fisherman/utils"
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

func (rule *ShellScript) GetPosition() byte {
	return Scripts
}

func (rule *ShellScript) Check(ctx internal.ExecutionContext, output io.Writer) error {
	return ctx.Shell().
		Exec(ctx, formatOutput(output, rule), rule.Shell, shell.ShScript{
			Commands: rule.Commands,
			Env:      rule.Env,
			Dir:      rule.Dir,
		})
}

func (rule *ShellScript) Compile(variables map[string]interface{}) {
	rule.BaseRule.Compile(variables)
	utils.FillTemplate(&rule.Dir, variables)
	utils.FillTemplate(&rule.Name, variables)
	utils.FillTemplatesArray(rule.Commands, variables)
	utils.FillTemplatesMap(rule.Env, variables)
}

func formatOutput(output io.Writer, rule *ShellScript) io.Writer {
	if rule.Output {
		return prefixwriter.New(output, fmt.Sprintf("[%s] ", rule.Name))
	}

	return ioutil.Discard
}
