package handlers_test

import (
	"fisherman/commands"
	"fisherman/config"
	"fisherman/handlers"
	iomock "fisherman/mocks/infrastructure"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrePushHandler(t *testing.T) {
	fakeRepository := iomock.Repository{}
	fakeRepository.On("GetCurrentBranch").Return("develop", nil)
	fakeRepository.On("GetLastTag").Return("0.0.0", nil)

	fakeShell := iomock.Shell{}

	assert.NotPanics(t, func() {
		err := handlers.PrePushHandler(&commands.CommandContext{
			Variables:  map[string]interface{}{},
			Config:     &config.HooksConfig{},
			Repository: &fakeRepository,
			Shell:      &fakeShell,
		}, []string{})
		assert.NoError(t, err)
	})
}
