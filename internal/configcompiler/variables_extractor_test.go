package configcompiler_test

import (
	"errors"
	hooks "fisherman/configuration"
	"fisherman/constants"
	"fisherman/infrastructure"
	"fisherman/internal/configcompiler"
	"fisherman/testing/mocks"
	"testing"

	"github.com/imdario/mergo"
	"github.com/stretchr/testify/assert"
)

type testSource struct {
	Data  string
	Error error
}

func TestVariables(t *testing.T) {
	defaultTagConfig := testSource{Data: "0.0.1"}
	defaultBranchConfig := testSource{Data: "refs/heads/master"}

	tests := []struct {
		name          string
		expectedError error
		fromBranch    string
		fromTag       string
		expectVars    map[string]interface{}
		tagConfig     testSource
		branchConfig  testSource
	}{
		{
			name:          "Load CurrentBranch from branch name",
			fromBranch:    "refs/heads/(?P<CurrentBranch>.*)",
			expectedError: nil,
			expectVars: map[string]interface{}{
				"CurrentBranch": "master",
				"UserName":      "username",
			},
			tagConfig:    defaultTagConfig,
			branchConfig: defaultBranchConfig,
		},
		{
			name:          "Load CurrentBranch and oweride Version from branch name",
			fromBranch:    "refs/(?P<FishermanVersion>.*)/(?P<CurrentBranch>.*)",
			expectedError: nil,
			expectVars: map[string]interface{}{
				"CurrentBranch":    "master",
				"FishermanVersion": "heads",
			},
			tagConfig:    defaultTagConfig,
			branchConfig: defaultBranchConfig,
		},
		{
			name:          "Load multiple variables from branch name",
			fromBranch:    "(?P<Ref>.*)/(?P<Head>.*)/(?P<CurrentBranch>.*)",
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
			name:          "Load variables from branch and tag names",
			fromBranch:    "refs/heads/(?P<CurrentBranch>.*)",
			fromTag:       "0.0.(?P<Test>.*)",
			expectedError: nil,
			expectVars: map[string]interface{}{
				"CurrentBranch": "master",
				"Test":          "1",
			},
			tagConfig:    defaultTagConfig,
			branchConfig: defaultBranchConfig,
		},
		{
			name:          "Load CurrentBranch from branch name and oweride Vertion from tag name",
			fromBranch:    "refs/heads/(?P<CurrentBranch>.*)",
			fromTag:       "(?P<FishermanVersion>.*)",
			expectedError: nil,
			expectVars: map[string]interface{}{
				"CurrentBranch":    "master",
				"FishermanVersion": "0.0.1",
			},
			tagConfig:    defaultTagConfig,
			branchConfig: defaultBranchConfig,
		},
		{
			name:          "Abort loading when branch source return error",
			fromBranch:    "refs/heads/(?P<CurrentBranch>.*)",
			fromTag:       "(?P<Version>.*)",
			expectedError: errors.New("Test"),
			expectVars: map[string]interface{}{
				"CurrentBranch": "master",
				"Version":       "0.0.1",
			},
			tagConfig: defaultTagConfig,
			branchConfig: testSource{
				Error: errors.New("Test"),
			},
		},
		{
			name:          "Abort loading when source return error",
			fromTag:       "(?P<Version>.*)",
			expectedError: errors.New("Test"),
			expectVars: map[string]interface{}{
				"CurrentBranch": "master",
				"Version":       "0.0.1",
			},
			tagConfig: testSource{
				Error: errors.New("Test"),
			},
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
			globals := map[string]interface{}{"GlobalVar": "global"}
			extractor := configcompiler.NewConfigExtractor(
				mocks.NewRepositoryMock(t).
					GetCurrentBranchMock.Return(tt.branchConfig.Data, tt.branchConfig.Error).
					GetLastTagMock.Return(tt.tagConfig.Data, tt.tagConfig.Error).
					GetUserMock.Return(infrastructure.User{Email: "email@test.com", UserName: "username"}, nil),
				globals,
				"demo",
			)

			vars, err := extractor.Variables(hooks.VariablesConfig{
				FromBranch:  tt.fromBranch,
				FromLastTag: tt.fromTag,
			})

			withGlobal(&tt.expectVars)

			assert.Equal(t, tt.expectedError, err)
			if tt.expectedError == nil {
				assert.EqualValues(t, tt.expectVars, vars)
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
