package rules

import (
	"context"
	"fisherman/internal/utils"
	"fisherman/pkg/shell/helpers"
	"io"
	"os/exec"

	"github.com/hashicorp/go-multierror"
	"golang.org/x/text/transform"
)

var CommandContext = exec.CommandContext

const ExecType = "exec"

type Exec struct {
	BaseRule `yaml:",inline"`
	Name     string            `yaml:"name"`
	Env      map[string]string `yaml:"env"`
	Output   bool              `yaml:"output"`
	Dir      string            `yaml:"dir"`
	Commands []CommandDef      `yaml:"commands"`
}

// TODO: Add custom YAML unmarshalling.
type CommandDef struct {
	Program  string            `yaml:"program"`
	Args     []string          `yaml:"args"`
	Env      map[string]string `yaml:"env"`
	Output   bool              `yaml:"output"`
	Encoding string            `yaml:"encoding"`
	Dir      string            `yaml:"dir"`
}

func (command *CommandDef) Compile(variables map[string]interface{}) {
	utils.FillTemplate(&command.Program, variables)
	utils.FillTemplatesArray(command.Args, variables)
	utils.FillTemplatesMap(command.Env, variables)
	utils.FillTemplate(&command.Dir, variables)
}

func (rule *Exec) GetPosition() byte {
	return Scripts
}

func (rule *Exec) GetPrefix() string {
	if utils.IsEmpty(rule.Name) {
		return ExecType
	}

	return rule.Name
}

func (rule *Exec) Check(ctx context.Context, output io.Writer) error {
	env := helpers.MergeEnv(rule.BaseRule.env, rule.Env)

	var resultError *multierror.Error
	for _, commandDef := range rule.Commands {
		encoding, err := getEncoding(commandDef.Encoding)
		if err != nil {
			return err
		}

		command := CommandContext(ctx, commandDef.Program, commandDef.Args...)
		command.Env = helpers.MergeEnv(env, commandDef.Env)
		command.Dir = utils.FirstNotEmpty(commandDef.Dir, rule.Dir, rule.BaseRule.cwd)
		command.Stdout = transform.NewWriter(output, encoding.NewDecoder())
		command.Stderr = transform.NewWriter(output, encoding.NewDecoder())

		if err := command.Run(); err != nil {
			resultError = multierror.Append(resultError, err)
		}
	}

	return resultError.ErrorOrNil()
}

func (rule *Exec) Compile(variables map[string]interface{}) {
	rule.BaseRule.Compile(variables)
	utils.FillTemplate(&rule.Dir, variables)
	utils.FillTemplate(&rule.Name, variables)
	utils.FillTemplatesMap(rule.Env, variables)
	for i := 0; i < len(rule.Commands); i++ {
		rule.Commands[i].Compile(variables)
	}
}
