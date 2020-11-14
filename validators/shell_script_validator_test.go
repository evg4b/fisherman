package validators_test

import (
	"context"
	"fisherman/infrastructure/shell"
	"fisherman/mocks"
	"fisherman/validators"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScriptValidator(t *testing.T) {
	script := shell.ScriptConfig{
		Name:     "test",
		Commands: []string{"command1", "command2"},
		Env: map[string]string{
			"demo":  "test1",
			"demo2": "test2",
		},
		Output: true,
	}

	expectedResult := shell.ExecResult{
		Name:     "test",
		Error:    nil,
		ExitCode: 9,
		Time:     time.Hour,
	}

	sh := mocks.NewShellMock(t).ExecMock.Inspect(func(ctx context.Context, shScript shell.ScriptConfig) {
		assert.NotNil(t, ctx)
		assert.Equal(t, script.Name, shScript.Name)
		assert.EqualValues(t, script.Commands, shScript.Commands)
		assert.EqualValues(t, script.Env, shScript.Env)
		assert.Equal(t, script.Output, shScript.Output)
	}).Return(expectedResult)

	ctx := mocks.NewAsyncContextMock(t).
		ShellMock.Return(sh)

	result := validators.ScriptValidator(ctx, script)

	assert.Equal(t, expectedResult.Name, result.Name)
	assert.Equal(t, expectedResult.Error, result.Error)
	assert.Equal(t, expectedResult.Time, result.Time)
}
