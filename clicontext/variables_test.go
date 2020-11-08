package clicontext_test

import (
	"errors"
	"fisherman/clicontext"
	"fisherman/config"
	"fisherman/constants"
	"fisherman/infrastructure"
	hooks_mock "fisherman/mocks/config/hooks"
	inf_mocks "fisherman/mocks/infrastructure"
	"testing"

	"github.com/imdario/mergo"
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
		expectedError   error
		fromBranch      string
		fromBranchVars  map[string]interface{}
		fromBranchError error
		fromTag         string
		fromTagVars     map[string]interface{}
		fromTagError    error
		expectVars      map[string]interface{}
		tagConfig       testSource
		branchConfig    testSource
	}{
		{
			name:       "Load CurrentBranch from branch name",
			fromBranch: "refs/heads/(?P<CurrentBranch>.*)",
			fromBranchVars: map[string]interface{}{
				"CurrentBranch": "master",
			},
			expectedError: nil,
			expectVars: map[string]interface{}{
				"CurrentBranch": "master",
			},
			tagConfig:    defaultTagConfig,
			branchConfig: defaultBranchConfig,
		},
		{
			name:       "Load CurrentBranch and oweride Version from branch name",
			fromBranch: "refs/(?P<Version>.*)/(?P<CurrentBranch>.*)",
			fromBranchVars: map[string]interface{}{
				"CurrentBranch": "master",
				"Version":       "heads",
			},
			expectedError: nil,
			expectVars: map[string]interface{}{
				"CurrentBranch": "master",
				"Version":       "heads",
			},
			tagConfig:    defaultTagConfig,
			branchConfig: defaultBranchConfig,
		},
		{
			name:       "Load multiple variables from branch name",
			fromBranch: "(?P<Ref>.*)/(?P<Head>.*)/(?P<CurrentBranch>.*)",
			fromBranchVars: map[string]interface{}{
				"CurrentBranch": "master",
				"Head":          "heads",
				"Ref":           "refs",
			},
			expectedError: nil,
			expectVars: map[string]interface{}{
				"CurrentBranch": "master",
				"Head":          "heads",
				"Ref":           "refs",
			},
			tagConfig:    defaultTagConfig,
			branchConfig: defaultBranchConfig,
		},
		{
			name:       "Load variables from branch and tag names",
			fromBranch: "refs/heads/(?P<CurrentBranch>.*)",
			fromBranchVars: map[string]interface{}{
				"CurrentBranch": "master",
			},
			fromTag: "0.0.(?P<Test>.*)",
			fromTagVars: map[string]interface{}{
				"Test": "1",
			},
			expectedError: nil,
			expectVars: map[string]interface{}{
				"CurrentBranch": "master",
				"Test":          "1",
			},
			tagConfig:    defaultTagConfig,
			branchConfig: defaultBranchConfig,
		},
		{
			name:       "Load CurrentBranch from branch name and oweride Vertion from tag name",
			fromBranch: "refs/heads/(?P<CurrentBranch>.*)",
			fromBranchVars: map[string]interface{}{
				"CurrentBranch": "master",
			},
			fromTag: "(?P<Version>.*)",
			fromTagVars: map[string]interface{}{
				"Version": "0.0.1",
			},
			expectedError: nil,
			expectVars: map[string]interface{}{
				"CurrentBranch": "master",
				"Version":       "0.0.1",
			},
			tagConfig:    defaultTagConfig,
			branchConfig: defaultBranchConfig,
		},
		{
			name:       "Abort loading when source return error",
			fromBranch: "refs/heads/(?P<CurrentBranch>.*)",
			fromBranchVars: map[string]interface{}{
				"CurrentBranch": "master",
			},
			fromBranchError: errors.New("Test"),
			fromTag:         "(?P<Version>.*)",
			fromTagVars: map[string]interface{}{
				"Version": "0.0.1",
			},
			expectedError: errors.New("Test"),
			expectVars: map[string]interface{}{
				"CurrentBranch": "master",
				"Version":       "0.0.1",
			},
			tagConfig:    defaultTagConfig,
			branchConfig: defaultBranchConfig,
		},
		{
			name:          "Abort loading when source return error",
			fromTag:       "(?P<Version>.*)",
			fromTagVars:   map[string]interface{}{},
			fromTagError:  errors.New("Test"),
			expectedError: errors.New("Test"),
			expectVars: map[string]interface{}{
				"CurrentBranch": "master",
				"Version":       "0.0.1",
			},
			tagConfig:    defaultTagConfig,
			branchConfig: defaultBranchConfig,
		},
		{
			name:          "Skip loading",
			expectedError: nil,
			expectVars:    map[string]interface{}{},
			tagConfig:     defaultTagConfig,
			branchConfig:  defaultBranchConfig,
		},
		{
			name:          "GetLastTag Error",
			expectedError: errors.New("GetTagError"),
			tagConfig: testSource{
				Error: errors.New("GetTagError"),
			},
			branchConfig: defaultBranchConfig,
		},
		{
			name:          "GetBranch Error",
			expectedError: errors.New("GetBranchError"),
			tagConfig:     defaultTagConfig,
			branchConfig: testSource{
				Error: errors.New("GetBranchError"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := inf_mocks.Repository{}
			repo.On("GetCurrentBranch").Return(tt.branchConfig.Data, tt.branchConfig.Error)
			repo.On("GetLastTag").Return(tt.tagConfig.Data, tt.tagConfig.Error)
			repo.On("GetUser").Return(infrastructure.User{Email: "email@test.com", UserName: "username"}, nil)

			vars := hooks_mock.VariablesExtractor{}
			vars.On("GetFromBranch", tt.branchConfig.Data).Return(tt.fromBranchVars, tt.fromBranchError)
			vars.On("GetFromTag", tt.tagConfig.Data).Return(tt.fromTagVars, tt.fromTagError)

			withGlobal(&tt.expectVars)

			ctx := clicontext.NewContext(clicontext.Args{
				Repository:      &repo,
				GlobalVariables: map[string]interface{}{"GlobalVar": "global"},
				Config:          &config.DefaultConfig,
				App: &clicontext.AppInfo{
					Cwd: "demo",
				},
			})

			err := ctx.LoadAdditionalVariables(&vars)

			assert.Equal(t, tt.expectedError, err)
			if tt.expectedError == nil {
				assert.EqualValues(t, tt.expectVars, ctx.Variables())
			}
		})
	}
}

func withGlobal(vars *map[string]interface{}) {
	err := mergo.Merge(vars, map[string]interface{}{
		constants.FishermanVersionVariable: "x.x.x",
		constants.CwdVariable:              "demo",
		constants.EmailVariable:            "email@test.com",
		constants.UserNameVariable:         "username",
		"GlobalVar":                        "global",
	})

	if err != nil {
		panic(err)
	}
}
