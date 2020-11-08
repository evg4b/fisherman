package handlers

import (
	"fisherman/clicontext"
	"fisherman/config"
	inf_mock "fisherman/mocks/infrastructure"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareCommitMsgHandler(t *testing.T) {
	fakeRepository := inf_mock.Repository{}
	fakeRepository.On("GetCurrentBranch").Return("develop", nil)
	fakeRepository.On("GetLastTag").Return("0.0.0", nil)

	fakeFS := inf_mock.FileSystem{}
	fakeFS.On("Read", ".git/MESSAGE").Return("[fisherman] test commit", nil)

	tests := []struct {
		name string
		args []string
		err  error
	}{
		{name: "base test", args: []string{}, err: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := clicontext.NewContext(clicontext.Args{
				Config:     &config.DefaultConfig,
				Repository: &fakeRepository,
				FileSystem: &fakeFS,
			})
			err := PrepareCommitMsgHandler(ctx, tt.args)
			assert.Equal(t, tt.err, err)
		})
	}
}
