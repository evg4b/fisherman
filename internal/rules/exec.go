package rules

import (
	"context"
	"github.com/evg4b/fisherman/internal/utils"
	"github.com/evg4b/fisherman/pkg/shell/helpers"
	"io"
	"os/exec"

	"github.com/hashicorp/go-multierror"
	"github.com/kballard/go-shellquote"
	"golang.org/x/text/transform"
	"gopkg.in/yaml.v3"
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

type CommandDef struct {
	Program  string            `yaml:"program"`
	Args     []string          `yaml:"args"`
	Env      map[string]string `yaml:"env"`
	Output   bool              `yaml:"output"`
	Encoding string            `yaml:"encoding"`
	Dir      string            `yaml:"dir"`
}

func (c *CommandDef) UnmarshalYAML(value *yaml.Node) error {
	*c = CommandDef{}

	var shortForm string
	if err := value.Decode(&shortForm); err == nil {
		parts, err := shellquote.Split(shortForm)
		if err != nil {
			return err
		}

		c.Program = parts[0]
		c.Args = parts[1:]

		return nil
	}

	type plain CommandDef

	return value.Decode((*plain)(c))
}

func (c *CommandDef) Compile(variables map[string]any) {
	utils.FillTemplate(&c.Program, variables)
	utils.FillTemplatesArray(c.Args, variables)
	utils.FillTemplatesMap(c.Env, variables)
	utils.FillTemplate(&c.Dir, variables)
}

func (r *Exec) GetPosition() byte {
	return Scripts
}

func (r *Exec) GetPrefix() string {
	if utils.IsEmpty(r.Name) {
		return ExecType
	}

	return r.Name
}

func (r *Exec) Check(ctx context.Context, output io.Writer) error {
	env := helpers.MergeEnv(r.BaseRule.env, r.Env)

	var resultError *multierror.Error
	for _, commandDef := range r.Commands {
		encoding, err := getEncoding(commandDef.Encoding)
		if err != nil {
			return err
		}

		command := CommandContext(ctx, commandDef.Program, commandDef.Args...)
		command.Env = helpers.MergeEnv(env, commandDef.Env)
		command.Dir = utils.FirstNotEmpty(commandDef.Dir, r.Dir, r.BaseRule.cwd)
		command.Stdout = transform.NewWriter(output, encoding.NewDecoder())
		command.Stderr = transform.NewWriter(output, encoding.NewDecoder())

		if err := command.Run(); err != nil {
			resultError = multierror.Append(resultError, err)
		}
	}

	return resultError.ErrorOrNil()
}

func (r *Exec) Compile(variables map[string]any) {
	r.BaseRule.Compile(variables)
	utils.FillTemplate(&r.Dir, variables)
	utils.FillTemplate(&r.Name, variables)
	utils.FillTemplatesMap(r.Env, variables)
	for i := 0; i < len(r.Commands); i++ {
		r.Commands[i].Compile(variables)
	}
}
