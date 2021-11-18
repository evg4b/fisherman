package rules

import (
	"fisherman/internal"
	"fisherman/internal/utils"
	"fisherman/pkg/shell"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"reflect"
	"runtime"

	"github.com/go-errors/errors"
	"gopkg.in/yaml.v3"
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

// TODO: Remove OS related overload struct. Use condition to run scripts in some OS.
type osRelatedShellScript struct {
	BaseRule `yaml:",inline"`
	Windows  BaseShell `yaml:"windows"`
	Linux    BaseShell `yaml:"linux"`
	Darwin   BaseShell `yaml:"macos"`
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

	shell := ctx.Shell()
	err := shell.Exec(ctx, formatOutput(output, rule), rule.Shell, script)
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

func (rule *ShellScript) UnmarshalYAML(value *yaml.Node) error {
	type plain ShellScript
	err := value.Decode((*plain)(rule))
	// FIXME: Fix problem with KnownFields parameter https://github.com/go-yaml/yaml/issues/460
	// DeepEqual checking needed to fix this issue.
	if err == nil && !reflect.DeepEqual(rule.BaseShell, BaseShell{}) {
		return nil
	}

	var osRelated osRelatedShellScript
	err = value.Decode(&osRelated)
	if err == nil {
		return fill(rule, &osRelated)
	}

	return err
}

func fill(rule *ShellScript, osRelated *osRelatedShellScript) error {
	switch runtime.GOOS {
	case "windows":
		rule.BaseRule = osRelated.BaseRule
		rule.BaseShell = osRelated.Windows
	case "linux":
		rule.BaseRule = osRelated.BaseRule
		rule.BaseShell = osRelated.Linux
	case "darwin":
		rule.BaseRule = osRelated.BaseRule
		rule.BaseShell = osRelated.Darwin
	default:
		return fmt.Errorf("system %s is not supported", runtime.GOOS)
	}

	return nil
}
