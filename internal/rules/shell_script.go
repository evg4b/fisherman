package rules

import (
	"fisherman/internal"
	"fisherman/internal/utils"
	"fisherman/pkg/shell"
	pkgutils "fisherman/pkg/utils"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"

	"github.com/go-errors/errors"
)

const ShellScriptType = "shell-script"

type BaseShell struct {
	Name     string            `yaml:"name"`
	Shell    string            `yaml:"shell"`
	Commands []string          `yaml:"commands"`
	Env      map[string]string `yaml:"env"`
	Output   bool              `yaml:"output"`
	Dir      string            `yaml:"dir"`
}

type ShellScript struct {
	BaseRule  `yaml:",inline"`
	BaseShell `yaml:",inline"`
}

func (rule *ShellScript) GetPosition() byte {
	return Scripts
}

func (rule *ShellScript) GetPrefix() string {
	return utils.FirstNotEmpty(rule.Name, rule.Type)
}

func (rule *ShellScript) Check(ctx internal.ExecutionContext, output io.Writer) error {
	formatterOutput := formatOutput(output, rule)
	host := shell.NewHost(
		ctx,
		shell.Cmd(), // TODO: Create factory method to resolve needed shell
		shell.WithEnv(pkgutils.MergeEnv(ctx.Env(), rule.Env)),
		shell.WithStdout(formatterOutput),
		shell.WithCwd(rule.Dir),
	)

	for _, command := range rule.Commands {
		_, err := fmt.Fprintln(host, command)
		if err != nil {
			return errors.Errorf("failed to exec shell script: %w", err)
		}
	}

	err := host.Wait()
	if err != nil {
		var exitCodeError *exec.ExitError
		if errors.As(err, &exitCodeError) {
			return rule.errorf("script finished with exit code %d", exitCodeError.ExitCode())
		}

		return errors.Errorf("failed to exec shell script: %w", err)
	}

	return nil
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
