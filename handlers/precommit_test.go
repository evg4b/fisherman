package handlers_test

import (
	"fisherman/clicontext"
	"fisherman/config"
	"fisherman/handlers"
	"fisherman/infrastructure"
	iomock "fisherman/mocks/infrastructure"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreCommitHandler(t *testing.T) {
	fakeRepository := iomock.Repository{}
	fakeRepository.On("GetCurrentBranch").Return("develop", nil)
	fakeRepository.On("GetLastTag").Return("0.0.0", nil)
	fakeRepository.On("GetUser").Return(infrastructure.User{}, nil)

	fakeShell := iomock.Shell{}

	assert.NotPanics(t, func() {
		err := handlers.PreCommitHandler(clicontext.NewContext(clicontext.Args{
			GlobalVariables: map[string]interface{}{},
			Config: &config.FishermanConfig{
				Hooks: config.HooksConfig{},
			},
			Repository: &fakeRepository,
			Shell:      &fakeShell,
			App:        &clicontext.AppInfo{},
		}), []string{})
		assert.NoError(t, err)
	})
}
