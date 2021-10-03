package shell_test

import (
	"context"
	"fisherman/pkg/shell"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScript_GetDuration(t *testing.T) {
	var expectedDuration time.Duration
	sh := shell.NewShell()

	script := shell.NewScript([]string{"echo 'demo'"})

	_ = sh.Exec(context.TODO(), io.Discard, shell.PlatformDefaultShell, script)

	duration := script.GetDuration()

	assert.Greater(t, duration, expectedDuration)
}

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
