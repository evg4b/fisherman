package hookfactory

import (
	"context"
	"errors"
	hooks "fisherman/configuration"
	"fisherman/infrastructure/shell"
	"fisherman/internal"
	"fisherman/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_stringWrapper(t *testing.T) {
	ctx := mocks.NewSyncContextMock(t)
	configuration := "configString"
	expectedError := errors.New("Error")

	wrapper := stringWrapper(func(actualCtx internal.SyncContext, actualString string) error {
		assert.Equal(t, ctx, actualCtx)
		assert.Equal(t, actualString, configuration)

		return expectedError
	}, configuration)

	err := wrapper(ctx)

	assert.Equal(t, expectedError, err)
}

func Test_boolWrapper(t *testing.T) {
	ctx := mocks.NewSyncContextMock(t)
	configuration := true
	expectedError := errors.New("Error")

	wrapper := boolWrapper(func(actualCtx internal.SyncContext, actualConfig bool) error {
		assert.Equal(t, ctx, actualCtx)
		assert.Equal(t, actualConfig, configuration)

		return expectedError
	}, configuration)

	err := wrapper(ctx)

	assert.Equal(t, expectedError, err)
}

func Test_scriptWrapper(t *testing.T) {
	ctx := mocks.NewAsyncContextMock(t)
	scriptName := "test-script"
	script := hooks.ScriptConfig{
		Commands: []string{"1", "2"},
		Env: map[string]string{
			"var": "value",
		},
		Output: true,
		Shell:  "test",
	}

	result := shell.ExecResult{
		Name:  "other-name",
		Error: errors.New("test-error"),
		Time:  time.Hour,
	}

	sh := mocks.NewShellMock(t).
		ExecMock.Inspect(func(actualCtx context.Context, bin string, actualScript shell.ShScriptConfig) {
		assert.Equal(t, ctx, actualCtx)
		assert.Equal(t, "test", bin)
		assert.Equal(t, scriptName, actualScript.Name)
		assert.ObjectsAreEqual(script, actualScript)
	}).Return(result)

	ctx.ShellMock.Return(sh)

	wrappers := scriptWrapper(hooks.ScriptsConfig{scriptName: script}, mocks.NewEngineMock(t))

	for _, wrapper := range wrappers {
		actualResult := wrapper(ctx)
		assert.ObjectsAreEqual(result, actualResult)
	}
}
