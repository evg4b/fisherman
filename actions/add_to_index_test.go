package actions_test

import (
	"errors"
	"fisherman/actions"
	"fisherman/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddToIndex_NotConfigured(t *testing.T) {
	next, err := actions.AddToIndex(mocks.NewSyncContextMock(t), []string{})

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
	})

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
	})

	assert.Error(t, err, "testError")
	assert.False(t, next)
}
