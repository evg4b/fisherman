package rules

import (
	"context"
	"fisherman/internal/utils"
	"fisherman/pkg/shell"
	"fisherman/pkg/shell/helpers"
	"fmt"
	"io"
	"os/exec"

	"github.com/go-errors/errors"
)

const ShellScriptType = "shell-script"

type ShellScript struct {
	BaseRule `yaml:",inline"`
	Name     string            `yaml:"name"`
	Shell    string            `yaml:"shell"`
	Commands []string          `yaml:"commands"`
	Env      map[string]string `yaml:"env"`
	Encoding string            `yaml:"encoding"`
	Output   bool              `yaml:"output"`
	Dir      string            `yaml:"dir"`
}

func (rule *ShellScript) GetPosition() byte {
	return Scripts
}

func (rule *ShellScript) GetPrefix() string {
	return utils.FirstNotEmpty(rule.Name, rule.Type)
}

func (rule *ShellScript) Check(ctx context.Context, output io.Writer) error {
	formatterOutput := formatOutput(output, rule)
	env := helpers.MergeEnv(rule.BaseRule.env, rule.Env)
	strategy, err := getShellStrategy(rule.Shell)
	if err != nil {
		return errors.Errorf("failed to cheate shell host: %w", err)
	}

	encoding, err := getEncoding(rule.Encoding)
	if err != nil {
		return errors.Errorf("failed to cheate shell host: %w", err)
	}

	host := shell.NewHost(
		ctx,
		strategy,
		shell.WithEnv(env),
		shell.WithStdout(formatterOutput),
		shell.WithCwd(rule.Dir),
		shell.WithEncoding(encoding),
	)

	for _, command := range rule.Commands {
		_, err := fmt.Fprintln(host, command)
		if err != nil {
			return errors.Errorf("failed to exec shell script: %w", err)
		}
	}

	err = host.Wait()
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

	return io.Discard
}

func getShellStrategy(name string) (shell.ShellStrategy, error) {
	if utils.IsEmpty(name) {
		return shell.Default(), nil
	}

	switch name {
	case "cmd":
		return shell.Cmd(), nil
	case "powershell":
		return shell.PowerShell(), nil
	case "bash":
		return shell.PowerShell(), nil
	default:
		return nil, errors.New("unsupported shell")
	}
}
