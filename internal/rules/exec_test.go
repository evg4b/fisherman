package rules_test

import (
	"fisherman/internal/rules"
	"fisherman/testing/mocks"
	"fisherman/testing/testutils"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExec_GetPosition(t *testing.T) {
	rule := rules.Exec{
		BaseRule: rules.BaseRule{Type: rules.ExecType},
	}

	actual := rule.GetPosition()

	assert.Equal(t, actual, rules.Scripts)
}

func TestExec_GetPrefix(t *testing.T) {
	tests := []struct {
		name     string
		ruleName string
		expected string
	}{
		{name: "user defined name", ruleName: "Prefix", expected: "Prefix"},
		{name: "default prefix", expected: rules.ExecType},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := rules.Exec{
				BaseRule: rules.BaseRule{Type: rules.ExecType},
				Name:     tt.ruleName,
			}

			actual := rule.GetPrefix()

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestExec_Check(t *testing.T) {
	fakeCommandContext, envWrapper := testutils.ConfigureFakeExec("TestExec_CheckHelper")

	rules.CommandContext = fakeCommandContext
	defer func() { rules.CommandContext = exec.CommandContext }()

	tests := []struct {
		name     string
		expected string
		commands []rules.CommandDef
	}{
		{
			name: "successfully command execution",
			commands: []rules.CommandDef{
				{Program: "go", Args: []string{"test", "./valid"}},
			},
		},
		{
			name:     "command finished with code 2",
			expected: "1 error occurred:\n\t* exit status 2\n\n",
			commands: []rules.CommandDef{
				{Program: "go", Args: []string{"test", "./..."}},
			},
		},
		{
			name: "successfully finished list of commands",
			commands: []rules.CommandDef{
				{Program: "go", Args: []string{"test", "./valid"}},
				{Program: "go", Args: []string{"test", "./another-valid"}},
			},
		},
		{
			name:     "failed one command from list",
			expected: "1 error occurred:\n\t* exit status 2\n\n",
			commands: []rules.CommandDef{
				{Program: "go", Args: []string{"test", "./..."}},
				{Program: "go", Args: []string{"test", "./valid"}},
				{Program: "go", Args: []string{"test", "./another-valid"}},
			},
		},
		{
			name:     "failed two command from list",
			expected: "2 errors occurred:\n\t* exit status 2\n\t* exit status 33\n\n",
			commands: []rules.CommandDef{
				{Program: "go", Args: []string{"test", "./..."}},
				{Program: "make", Args: []string{"build"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := rules.Exec{
				BaseRule: rules.BaseRule{Type: rules.ExecType},
				Commands: tt.commands,
			}

			ctx := mocks.NewExecutionContextMock(t).
				EnvMock.Return(envWrapper([]string{})).
				CwdMock.Return("/")

			actual := rule.Check(ctx, io.Discard)

			testutils.CheckError(t, tt.expected, actual)
		})
	}
}

func TestExec_CheckHelper(t *testing.T) {
	testutils.ExecTestHandler(t, map[string]func(){
		"go test ./...":           func() { os.Exit(2) },
		"go test ./valid":         func() { os.Exit(0) },
		"go test ./another-valid": func() { os.Exit(0) },
		"make build":              func() { os.Exit(33) },
	})
}
