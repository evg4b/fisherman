package actions_test

import (
	"errors"
	"fisherman/actions"
	"fisherman/mocks"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
)

func TestAddToIndex_NotConfigured(t *testing.T) {
	next, err := actions.AddToIndex(mocks.NewSyncContextMock(t), []string{}, false)

	assert.NoError(t, err)
	assert.True(t, next)
}

func TestAddToIndex_CorrectAddToIndex(t *testing.T) {
	repo := mocks.NewRepositoryMock(t).
		AddGlobMock.When("glob1/*.go").Then(nil).
		AddGlobMock.When("*.css").Then(nil).
		AddGlobMock.When("mocks").Then(nil)

	ctx := mocks.NewSyncContextMock(t).RepositoryMock.Return(repo)

	next, err := actions.AddToIndex(ctx, []string{
		"glob1/*.go",
		"*.css",
		"mocks",
	}, false)

	assert.NoError(t, err)
	assert.True(t, next)
}

func TestAddToIndex_FailedAddToIndex(t *testing.T) {
	repo := mocks.NewRepositoryMock(t).
		AddGlobMock.When("glob1/*.go").Then(nil).
		AddGlobMock.When("*.css").Then(errors.New("testError")).
		AddGlobMock.When("mocks").Then(nil)

	ctx := mocks.NewSyncContextMock(t).RepositoryMock.Return(repo)

	next, err := actions.AddToIndex(ctx, []string{
		"glob1/*.go",
		"*.css",
		"mocks",
	}, false)

	assert.Error(t, err, "testError")
	assert.False(t, next)
}

func TestAddToIndex_FailedAddToIndexOptional(t *testing.T) {
	repo := mocks.NewRepositoryMock(t).
		AddGlobMock.When("glob1/*.go").Then(nil).
		AddGlobMock.When("*.css").Then(git.ErrGlobNoMatches).
		AddGlobMock.When("mocks").Then(nil)

	ctx := mocks.NewSyncContextMock(t).RepositoryMock.Return(repo)

	tests := []struct {
		name     string
		optional bool
	}{
		{name: "Optional true", optional: true},
		{name: "Optional false", optional: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			next, err := actions.AddToIndex(ctx, []string{
				"glob1/*.go",
				"*.css",
				"mocks",
			}, tt.optional)

			if tt.optional {
				assert.NoError(t, err)
				assert.True(t, next)
			} else {
				assert.Equal(t, err, git.ErrGlobNoMatches)
				assert.False(t, next)
			}
		})
	}

}
