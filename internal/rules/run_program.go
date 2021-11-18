package rules

import (
	"fisherman/internal"
	"fisherman/internal/utils"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

const RunProgramType = "run-program"

type RunProgram struct {
	BaseRule `yaml:",inline"`
	Name     string            `yaml:"name"`
	Program  string            `yaml:"program"`
	Args     []string          `yaml:"args"`
	Env      map[string]string `yaml:"env"`
	Output   bool              `yaml:"output"`
	Dir      string            `yaml:"dir"`
}

func (rule *RunProgram) GetPosition() byte {
	return Scripts
}

func (rule *RunProgram) GetPrefix() string {
	if !utils.IsEmpty(rule.Name) {
		return rule.Name
	}

	generatedPrefix := fmt.Sprintf("%s %s", rule.Program, strings.Join(rule.Args, " "))

	return normalizePrefix(generatedPrefix)
}

func (rule *RunProgram) Check(ctx internal.ExecutionContext, output io.Writer) error {
	envList := os.Environ()
	for key, value := range rule.Env {
		envList = append(envList, fmt.Sprintf("%s=%s", key, value))
	}

	command := exec.CommandContext(ctx, rule.Program, rule.Args...) // nolint gosec
	command.Env = envList
	command.Dir = utils.GetOrDefault(rule.Dir, ctx.Cwd())

	// TODO: Add custom encoding for different shell
	command.Stdout = output
	command.Stderr = output

	return command.Run() // TODO: Add duration and output
}

func (rule *RunProgram) Compile(variables map[string]interface{}) {
	rule.BaseRule.Compile(variables)
	utils.FillTemplate(&rule.Dir, variables)
	utils.FillTemplatesArray(rule.Args, variables)
	utils.FillTemplatesMap(rule.Env, variables)
}
