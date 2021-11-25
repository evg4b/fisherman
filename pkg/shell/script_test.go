package shell_test

import (
	"fisherman/pkg/shell"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScript_EnvironmentVariables(t *testing.T) {
	expectedVars := map[string]string{
		"test": "tests",
	}

	script := shell.NewScript([]string{}).
		SetEnvironmentVariables(expectedVars)

	actualVars := script.GetEnvironmentVariables()

	assert.Equal(t, expectedVars, actualVars)
}

func TestScript_Directory(t *testing.T) {
	expectedDirectory := "~/projects/fisherman"

	script := shell.NewScript([]string{}).
		SetDirectory(expectedDirectory)

	actualDirectory := script.GetDirectory()

	assert.Equal(t, expectedDirectory, actualDirectory)
}

func TestScript_GetCommands(t *testing.T) {
	expectedCommands := []string{
		"echo 'test string'",
	}

	script := shell.NewScript(expectedCommands)

	actualCommands := script.GetCommands()

	assert.Equal(t, expectedCommands, actualCommands)
}
