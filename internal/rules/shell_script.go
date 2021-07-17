package rules

import (
	"fisherman/internal"
	"fisherman/internal/utils"
	"fisherman/pkg/shell"
	"io"
	"io/ioutil"
)

const ShellScriptType = "shell-script"

type ShellScript struct {
	BaseRule `yaml:",inline"`
	Name     string            `yaml:"name"`
	Shell    string            `yaml:"shell"`
	Commands []string          `yaml:"commands"`
	Env      map[string]string `yaml:"env"`
	Output   bool              `yaml:"output"`
	Dir      string            `yaml:"dir"`
}

func (rule *ShellScript) GetPosition() byte {
	return Scripts
}

func (rule *ShellScript) GetPrefix() string {
	return utils.GetOrDefault(rule.Name, rule.Type)
}

func (rule *ShellScript) Check(ctx internal.ExecutionContext, output io.Writer) error {
	script := shell.NewScript(rule.Commands).
		SetEnvironmentVariables(rule.Env).
		SetDirectory(rule.Dir)

	return ctx.Shell().
		Exec(ctx, formatOutput(output, rule), rule.Shell, script)
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
		return output
	}

	return ioutil.Discard
}
