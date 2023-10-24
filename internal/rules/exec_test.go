package rules_test

import (
	"context"
	"fisherman/testing/testutils"
	"io"
	"os"
	"os/exec"
	"testing"

	. "fisherman/internal/rules"

	"github.com/stretchr/testify/assert"
)

func TestExec_GetPosition(t *testing.T) {
	rule := Exec{BaseRule: BaseRule{Type: ExecType}}

	actual := rule.GetPosition()

	assert.Equal(t, actual, Scripts)
}

func TestExec_GetPrefix(t *testing.T) {
	tests := []struct {
		name     string
		ruleName string
		expected string
	}{
		{name: "user defined name", ruleName: "Prefix", expected: "Prefix"},
		{name: "default prefix", expected: ExecType},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := Exec{
				BaseRule: BaseRule{Type: ExecType},
				Name:     tt.ruleName,
			}

			actual := rule.GetPrefix()

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestExec_Check(t *testing.T) {
	fakeCommandContext, envWrapper := testutils.ConfigureFakeExec("TestExec_CheckHelper")

	CommandContext = fakeCommandContext
	defer func() { CommandContext = exec.CommandContext }()

	tests := []struct {
		name        string
		commands    []CommandDef
		env         map[string]string
		expectedErr string
	}{
		{
			name: "unknown encoding",
			commands: []CommandDef{
				{Program: "go", Args: []string{"test", "./valid"}, Encoding: "unknown"},
			},
			expectedErr: "'unknown' is unknown encoding",
		},
		{
			name: "correct apply encoding by name",
			commands: []CommandDef{
				{Program: "go", Args: []string{"test", "./valid"}, Encoding: "cp862"},
			},
		},
		{
			name: "successfully command execution",
			commands: []CommandDef{
				{Program: "go", Args: []string{"test", "./valid"}},
			},
		},
		{
			name:        "command finished with code 2",
			expectedErr: "1 error occurred:\n\t* exit status 2\n\n",
			commands: []CommandDef{
				{Program: "go", Args: []string{"test", "./..."}},
			},
		},
		{
			name: "successfully finished list of commands",
			commands: []CommandDef{
				{Program: "go", Args: []string{"test", "./valid"}},
				{Program: "go", Args: []string{"test", "./another-valid"}},
			},
		},
		{
			name:        "failed one command from list",
			expectedErr: "1 error occurred:\n\t* exit status 2\n\n",
			commands: []CommandDef{
				{Program: "go", Args: []string{"test", "./..."}},
				{Program: "go", Args: []string{"test", "./valid"}},
				{Program: "go", Args: []string{"test", "./another-valid"}},
			},
		},
		{
			name:        "failed two command from list",
			expectedErr: "2 errors occurred:\n\t* exit status 2\n\t* exit status 33\n\n",
			commands: []CommandDef{
				{Program: "go", Args: []string{"test", "./..."}},
				{Program: "make", Args: []string{"build"}},
			},
		},
		{
			name: "successfully passed global env",
			commands: []CommandDef{
				{Program: "go", Args: []string{"test", "env-global"}},
			},
		},
		{
			name: "successfully passed rule env",
			env: map[string]string{
				"RULE_ENV": "This is rule env",
			},
			commands: []CommandDef{
				{Program: "go", Args: []string{"test", "env-rule"}},
			},
		},
		{
			name: "successfully passed rule env",
			env: map[string]string{
				"RULE_ENV": "This is rule env",
			},
			commands: []CommandDef{
				{Program: "go", Args: []string{"test", "env-rule"}},
			},
		},
		{
			name: "successfully passed command env",
			commands: []CommandDef{
				{
					Program: "go",
					Args:    []string{"test", "env-command"},
					Env: map[string]string{
						"COMMAND_ENV": "This is command env",
					},
				},
			},
		},
		{
			name: "successfully oweride env variables",
			env: map[string]string{
				"RULE_ENV": "This is rule env",
				"VAR_2":    "Rule value 2",
				"VAR_3":    "Rule value 3",
			},
			commands: []CommandDef{
				{
					Program: "go",
					Args:    []string{"test", "env-oweride"},
					Env: map[string]string{
						"VAR_3": "Command value 3",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := makeRule(
				&Exec{
					BaseRule: BaseRule{Type: ExecType},
					Commands: tt.commands,
					Env:      tt.env,
				},
				WithCwd("/"),
				WithEnv(envWrapper([]string{
					"GLOBAL_ENV=6482376487",
					"PATH=Overwritten variable",
					"VAR_1=Global value 1",
					"VAR_2=Global value 2",
					"VAR_3=Global value 3",
				})),
			)

			actual := rule.Check(context.TODO(), io.Discard)

			testutils.AssertError(t, tt.expectedErr, actual)
		})
	}
}

func TestExec_Compile(t *testing.T) {
	rule := Exec{
		BaseRule: BaseRule{
			Type:      ExecType,
			Condition: "{{VAR1}}=\"TEST\"",
		},
		Name: "{{VAR1}}-{{VAR2}}",
		Env: map[string]string{
			"VAR1": "[{{VAR1}}]",
		},
		Output: false,
		Dir:    "/user/{{VAR2}}/test",
		Commands: []CommandDef{
			{
				Program: "app-{{Version}}",
				Args:    []string{"build", "--value={{VAR2}}"},
				Env: map[string]string{
					"VAR1": "[{{VAR1}}]",
				},
				Output: false,
				Dir:    "/user/{{VAR1}}/{{Version}}",
			},
		},
	}

	rule.Compile(map[string]any{
		"VAR1":    "TEST",
		"VAR2":    "DEMO",
		"Version": "3-5-17",
	})

	assert.Equal(t, Exec{
		BaseRule: BaseRule{
			Type:      ExecType,
			Condition: "TEST=\"TEST\"",
		},
		Name: "TEST-DEMO",
		Env: map[string]string{
			"VAR1": "[TEST]",
		},
		Output: false,
		Dir:    "/user/DEMO/test",
		Commands: []CommandDef{
			{
				Program: "app-3-5-17",
				Args:    []string{"build", "--value=DEMO"},
				Env: map[string]string{
					"VAR1": "[TEST]",
				},
				Output: false,
				Dir:    "/user/TEST/3-5-17",
			},
		},
	}, rule)
}

func TestExec_CheckHelper(t *testing.T) {
	testutils.ExecTestHandler(t, map[string]func(){
		"go test ./...":           func() { os.Exit(2) },
		"go test ./valid":         func() { os.Exit(0) },
		"go test ./another-valid": func() { os.Exit(0) },
		"make build":              func() { os.Exit(33) },
		"go test env-global": func() {
			assert.Equal(t, "6482376487", os.Getenv("GLOBAL_ENV"))
			assert.Equal(t, "Overwritten variable", os.Getenv("PATH"))
		},
		"go test env-rule": func() {
			assert.Equal(t, "This is rule env", os.Getenv("RULE_ENV"))
		},
		"go test env-command": func() {
			assert.Equal(t, "This is command env", os.Getenv("COMMAND_ENV"))
		},
		"go test env-oweride": func() {
			assert.Equal(t, "Global value 1", os.Getenv("VAR_1"))
			assert.Equal(t, "Rule value 2", os.Getenv("VAR_2"))
			assert.Equal(t, "Command value 3", os.Getenv("VAR_3"))
		},
	})
}

func TestCommandDef_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expected    CommandDef
		expectedErr string
	}{
		{
			name: "full form",
			value: `
program: go
args: [ build, main.go ]
env:
  GOOS: linux
output: true
encoding: utf8
dir: '~'
`,
			expected: CommandDef{
				Program: "go",
				Args:    []string{"build", "main.go"},
				Env: map[string]string{
					"GOOS": "linux",
				},
				Output:   true,
				Encoding: "utf8",
				Dir:      "~",
			},
		},
		{
			name:  "one string form",
			value: "go build main.go",
			expected: CommandDef{
				Program: "go",
				Args:    []string{"build", "main.go"},
			},
		},
		{
			name:  "one string form with spaces in binary path",
			value: "\"'/usr/test user/go' build main.go\"",
			expected: CommandDef{
				Program: "/usr/test user/go",
				Args:    []string{"build", "main.go"},
			},
		},
		{
			name:  "one string form with spaces in args",
			value: "\"go build main.go -ldflags '-s -w'\"",
			expected: CommandDef{
				Program: "go",
				Args:    []string{"build", "main.go", "-ldflags", "-s -w"},
			},
		},
		{
			name:        "one string form with spaces in args",
			value:       "\"go build main.go -ldflags '-s -w\"",
			expectedErr: "Unterminated single-quoted string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var def CommandDef

			err := testutils.DecodeYaml(tt.value, &def)

			assert.Equal(t, tt.expected, def)
			testutils.AssertError(t, tt.expectedErr, err)
		})
	}
}
