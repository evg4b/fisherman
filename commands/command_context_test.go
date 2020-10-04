package commands_test

import (
	"fisherman/commands"
	"fisherman/config"
	"fisherman/config/hooks"
	mocks "fisherman/mocks/infrastructure"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandContext_LoadAdditionalVariables(t *testing.T) {
	repo := mocks.Repository{}
	repo.On("GetCurrentBranch").Return("refs/heads/master", nil)
	repo.On("GetLastTag").Return("0.0.1", nil)

	tests := []struct {
		name          string
		args          hooks.Variables
		expectedError error
		expectVars    map[string]interface{}
	}{
		{
			name: "asd",
			args: hooks.Variables{
				FromBranch: "refs/heads/(?P<CurrentBranch>.*)",
			},
			expectedError: nil,
			expectVars: map[string]interface{}{
				"CurrentBranch": "master",
				"Version":       "version",
			},
		},
		{
			name: "asd",
			args: hooks.Variables{
				FromBranch: "refs/(?P<Version>.*)/(?P<CurrentBranch>.*)",
			},
			expectedError: nil,
			expectVars: map[string]interface{}{
				"CurrentBranch": "master",
				"Version":       "heads",
			},
		},
		{
			name: "asd",
			args: hooks.Variables{
				FromBranch: "(?P<Ref>.*)/(?P<Head>.*)/(?P<CurrentBranch>.*)",
			},
			expectedError: nil,
			expectVars: map[string]interface{}{
				"CurrentBranch": "master",
				"Version":       "version",
				"Head":          "heads",
				"Ref":           "refs",
			},
		},
		{
			name: "asd",
			args: hooks.Variables{
				FromBranch:  "refs/heads/(?P<CurrentBranch>.*)",
				FromLastTag: "0.0.(?P<Test>.*)",
			},
			expectedError: nil,
			expectVars: map[string]interface{}{
				"CurrentBranch": "master",
				"Version":       "version",
				"Test":          "1",
			},
		},
		{
			name: "asd",
			args: hooks.Variables{
				FromBranch:  "refs/heads/(?P<CurrentBranch>.*)",
				FromLastTag: "(?P<Version>.*)",
			},
			expectedError: nil,
			expectVars: map[string]interface{}{
				"CurrentBranch": "master",
				"Version":       "0.0.1",
			},
		},
		{
			name: "asd",
			args: hooks.Variables{
				FromBranch:  "",
				FromLastTag: "",
			},
			expectedError: nil,
			expectVars: map[string]interface{}{
				"Version": "version",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := commands.NewContext(commands.CliCommandContextParams{
				Repository: &repo,
				Variables:  map[string]interface{}{"Version": "version"},
				Config:     &config.DefaultConfig,
			})

			err := ctx.LoadAdditionalVariables(&tt.args)
			assert.Equal(t, tt.expectedError, err)

			for expectedKey, expectedValue := range tt.expectVars {
				assert.Equal(t, ctx.Variables[expectedKey], expectedValue)
			}
		})
	}
}
