package clicontext_test

import (
	"errors"
	"fisherman/clicontext"
	"fisherman/config"
	"fisherman/config/hooks"
	mocks "fisherman/mocks/infrastructure"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testSource struct {
	Data  string
	Error error
}

func TestCommandContext_LoadAdditionalVariables(t *testing.T) {
	defaultTagConfig := testSource{Data: "0.0.1"}
	defaultBranchConfig := testSource{Data: "refs/heads/master"}

	tests := []struct {
		name            string
		variables       hooks.Variables
		expectedError   error
		expectVars      map[string]interface{}
		getTagConfig    testSource
		getBranchConfig testSource
	}{
		{
			name: "Load CurrentBranch from branch name",
			variables: hooks.Variables{
				FromBranch: "refs/heads/(?P<CurrentBranch>.*)",
			},
			expectedError: nil,
			expectVars: map[string]interface{}{
				"CurrentBranch": "master",
				"Version":       "version",
			},
			getTagConfig:    defaultTagConfig,
			getBranchConfig: defaultBranchConfig,
		},
		{
			name: "Load CurrentBranch and oweride Version from branch name",
			variables: hooks.Variables{
				FromBranch: "refs/(?P<Version>.*)/(?P<CurrentBranch>.*)",
			},
			expectedError: nil,
			expectVars: map[string]interface{}{
				"CurrentBranch": "master",
				"Version":       "heads",
			},
			getTagConfig:    defaultTagConfig,
			getBranchConfig: defaultBranchConfig,
		},
		{
			name: "Load multiple variables from branch name",
			variables: hooks.Variables{
				FromBranch: "(?P<Ref>.*)/(?P<Head>.*)/(?P<CurrentBranch>.*)",
			},
			expectedError: nil,
			expectVars: map[string]interface{}{
				"CurrentBranch": "master",
				"Version":       "version",
				"Head":          "heads",
				"Ref":           "refs",
			},
			getTagConfig:    defaultTagConfig,
			getBranchConfig: defaultBranchConfig,
		},
		{
			name: "Load variables from branch and tag names",
			variables: hooks.Variables{
				FromBranch:  "refs/heads/(?P<CurrentBranch>.*)",
				FromLastTag: "0.0.(?P<Test>.*)",
			},
			expectedError: nil,
			expectVars: map[string]interface{}{
				"CurrentBranch": "master",
				"Version":       "version",
				"Test":          "1",
			},
			getTagConfig:    defaultTagConfig,
			getBranchConfig: defaultBranchConfig,
		},
		{
			name: "Load CurrentBranch from branch name and oweride Vertion from tag name",
			variables: hooks.Variables{
				FromBranch:  "refs/heads/(?P<CurrentBranch>.*)",
				FromLastTag: "(?P<Version>.*)",
			},
			expectedError: nil,
			expectVars: map[string]interface{}{
				"CurrentBranch": "master",
				"Version":       "0.0.1",
			},
			getTagConfig:    defaultTagConfig,
			getBranchConfig: defaultBranchConfig,
		},
		{
			name: "Skip loading",
			variables: hooks.Variables{
				FromBranch:  "",
				FromLastTag: "",
			},
			expectedError: nil,
			expectVars: map[string]interface{}{
				"Version": "version",
			},
			getTagConfig:    defaultTagConfig,
			getBranchConfig: defaultBranchConfig,
		},
		{
			name: "GetLastTag Error",
			variables: hooks.Variables{
				FromBranch:  "",
				FromLastTag: "",
			},
			expectedError: errors.New("GetTagError"),
			getTagConfig: testSource{
				Error: errors.New("GetTagError"),
			},
			getBranchConfig: defaultBranchConfig,
		},
		{
			name: "GetBranch Error",
			variables: hooks.Variables{
				FromBranch:  "",
				FromLastTag: "",
			},
			expectedError: errors.New("GetBranchError"),
			getTagConfig:  defaultTagConfig,
			getBranchConfig: testSource{
				Error: errors.New("GetBranchError"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.Repository{}
			repo.On("GetCurrentBranch").Return(tt.getBranchConfig.Data, tt.getBranchConfig.Error)
			repo.On("GetLastTag").Return(tt.getTagConfig.Data, tt.getTagConfig.Error)

			ctx := clicontext.NewContext(clicontext.Args{
				Repository: &repo,
				Variables:  map[string]interface{}{"Version": "version"},
				Config:     &config.DefaultConfig,
			})

			err := ctx.LoadAdditionalVariables(&tt.variables)

			assert.Equal(t, tt.expectedError, err)
			if tt.expectedError == nil {
				assert.EqualValues(t, tt.expectVars, ctx.Variables)
			}
		})
	}
}
