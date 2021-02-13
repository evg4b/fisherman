package configuration_test

import (
	"errors"
	"fisherman/configuration"
	"fisherman/constants"
	"fisherman/infrastructure/log"
	"fisherman/internal/rules"
	"fisherman/testing/mocks"
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestConfigLoader_FindConfigFiles(t *testing.T) {
	usr := user.User{HomeDir: filepath.Join("/", "usr", "home")}
	cwd := filepath.Join("/", "usr", "home", "documents", "repo")

	localConfig := configuration.GetConfigFolder(&usr, cwd, configuration.LocalMode)
	repoConfig := configuration.GetConfigFolder(&usr, cwd, configuration.RepoMode)
	globalConfig := configuration.GetConfigFolder(&usr, cwd, configuration.GlobalMode)

	tests := []struct {
		name        string
		loader      *configuration.ConfigLoader
		expected    map[string]string
		expectedErr error
	}{
		{
			name: "no config files",
			loader: configuration.NewLoader(
				&usr,
				cwd,
				mocks.NewFileSystemMock(t).
					FindMock.When(localConfig, constants.AppConfigNames).Then([]string{}, nil).
					FindMock.When(repoConfig, constants.AppConfigNames).Then([]string{}, nil).
					FindMock.When(globalConfig, constants.AppConfigNames).Then([]string{}, nil)),
			expectedErr: errors.New("no configuration found"),
		},
		{
			name: "mere then one config file",
			loader: configuration.NewLoader(
				&usr,
				cwd,
				mocks.NewFileSystemMock(t).
					FindMock.When(localConfig, constants.AppConfigNames).Then([]string{"file1", "file2"}, nil).
					FindMock.When(repoConfig, constants.AppConfigNames).Then([]string{}, nil).
					FindMock.When(globalConfig, constants.AppConfigNames).Then([]string{}, nil)),
			expectedErr: fmt.Errorf("more then one config file specifies in folder '%s'", localConfig),
		},
		{
			name: "correct files loading",
			loader: configuration.NewLoader(
				&usr,
				cwd,
				mocks.NewFileSystemMock(t).
					FindMock.When(localConfig, constants.AppConfigNames).Then([]string{"file1"}, nil).
					FindMock.When(repoConfig, constants.AppConfigNames).Then([]string{"file2"}, nil).
					FindMock.When(globalConfig, constants.AppConfigNames).Then([]string{"file3"}, nil)),
			expected: map[string]string{
				configuration.LocalMode:  "file1",
				configuration.RepoMode:   "file2",
				configuration.GlobalMode: "file3",
			},
		},
		{
			name: "unknown I/O error",
			loader: configuration.NewLoader(
				&usr,
				cwd,
				mocks.NewFileSystemMock(t).
					FindMock.When(localConfig, constants.AppConfigNames).Then([]string{"file1"}, errors.New("Unknown error")).
					FindMock.When(repoConfig, constants.AppConfigNames).Then([]string{"file2"}, nil).
					FindMock.When(globalConfig, constants.AppConfigNames).Then([]string{"file3"}, nil)),
			expectedErr: errors.New("Unknown error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := tt.loader.FindConfigFiles()

			assert.Equal(t, tt.expectedErr, err)
			if tt.expectedErr == nil {
				assert.EqualValues(t, tt.expected, actual)
			}
		})
	}
}

func TestGetConfigFolder(t *testing.T) {
	usr := user.User{HomeDir: filepath.Join("/", "usr", "home")}
	cwd := filepath.Join("/", "usr", "home", "documents", "repo")

	tests := []struct {
		name        string
		usr         *user.User
		cwd         string
		mode        string
		expected    string
		shouldPanic bool
	}{
		{
			name:     "local mode",
			usr:      &usr,
			cwd:      cwd,
			mode:     configuration.LocalMode,
			expected: filepath.Join(cwd, ".git"),
		},
		{
			name:     "global mode",
			usr:      &usr,
			cwd:      cwd,
			mode:     configuration.GlobalMode,
			expected: usr.HomeDir,
		},
		{
			name:     "repository mode",
			usr:      &usr,
			cwd:      cwd,
			mode:     configuration.RepoMode,
			expected: cwd,
		},
		{
			name:        "unknown mode",
			usr:         &usr,
			cwd:         cwd,
			mode:        "unknown mode",
			shouldPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				assert.Panics(t, func() {
					_ = configuration.GetConfigFolder(tt.usr, tt.cwd, tt.mode)
				})
			} else {
				actual := configuration.GetConfigFolder(tt.usr, tt.cwd, tt.mode)

				assert.Equal(t, tt.expected, actual)
			}
		})
	}
}

func TestConfigLoader_Load(t *testing.T) {
	usr := user.User{HomeDir: filepath.Join("/", "usr", "home")}
	cwd := filepath.Join("/", "usr", "home", "documents", "repo")

	config := `
variables:
  name: value

hooks:
  pre-push:
    rules:
      - type: shell-script
        name: Demo
        commands:
         - echo '1213123' >> log.txt
         - exit 1`

	reader := ioutil.NopCloser(strings.NewReader(config))

	tests := []struct {
		name        string
		loader      *configuration.ConfigLoader
		files       map[string]string
		expected    *configuration.FishermanConfig
		expectedErr error
	}{
		{
			name:   "",
			loader: configuration.NewLoader(&usr, cwd, mocks.NewFileSystemMock(t)),
			files:  map[string]string{},
			expected: &configuration.FishermanConfig{
				Output: log.DefaultOutputConfig,
			},
		},
		{
			name: "file reader error",
			loader: configuration.NewLoader(
				&usr,
				cwd,
				mocks.NewFileSystemMock(t).ReaderMock.When("GlobalConfig").Then(reader, errors.New("error"))),
			files: map[string]string{
				configuration.GlobalMode: "GlobalConfig",
			},
			expectedErr: errors.New("error"),
		},
		{
			name: "correct decoding",
			loader: configuration.NewLoader(
				&usr,
				cwd,
				mocks.NewFileSystemMock(t).ReaderMock.When("GlobalConfig").Then(reader, nil)),
			files: map[string]string{
				configuration.GlobalMode: "GlobalConfig",
			},
			expected: &configuration.FishermanConfig{
				Output: log.DefaultOutputConfig,
				GlobalVariables: map[string]interface{}{
					"name": "value",
				},
				Hooks: configuration.HooksConfig{
					PrePushHook: &configuration.PrePushHookConfig{
						RulesSection: configuration.RulesSection{
							Rules: []configuration.Rule{
								&rules.ShellScript{
									BaseRule: rules.BaseRule{
										Type: "shell-script",
									},
									Name: "Demo",
									Commands: []string{
										"echo '1213123' >> log.txt",
										"exit 1",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "decoding error",
			loader: configuration.NewLoader(
				&usr,
				cwd,
				mocks.NewFileSystemMock(t).ReaderMock.When("GlobalConfig").Then(ioutil.NopCloser(strings.NewReader("asd['")), nil)),
			files: map[string]string{
				configuration.GlobalMode: "GlobalConfig",
			},
			expectedErr: &yaml.TypeError{
				Errors: []string{"line 1: cannot unmarshal !!str `asd['` into configuration.FishermanConfig"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := tt.loader.Load(tt.files)

			assert.Equal(t, tt.expectedErr, err)
			if tt.expectedErr == nil {
				assert.Equal(t, tt.expected, actual)
			}
		})
	}
}

func TestConfigLoader_Load_Correct_Merging(t *testing.T) {
	usr := user.User{HomeDir: filepath.Join("/", "usr", "home")}
	cwd := filepath.Join("/", "usr", "home", "documents", "repo")

	fs := mocks.NewFileSystemMock(t)
	fs.ReaderMock.When("global.yaml").Then(ioutil.NopCloser(strings.NewReader(`
variables:
  var1: global
  var2: global
  var3: global
  var4: global`)), nil)
	fs.ReaderMock.When("repo.yaml").Then(ioutil.NopCloser(strings.NewReader(`
variables:
  var1: repo
  var2: repo`)), nil)
	fs.ReaderMock.When("local.yaml").Then(ioutil.NopCloser(strings.NewReader(`
variables:
  var1: local
  var3: local
  `)), nil)

	loader := configuration.NewLoader(&usr, cwd, fs)

	files := map[string]string{
		configuration.GlobalMode: "global.yaml",
		configuration.LocalMode:  "local.yaml",
		configuration.RepoMode:   "repo.yaml",
	}

	actual, err := loader.Load(files)

	assert.NoError(t, err)
	assert.Equal(t, &configuration.FishermanConfig{
		Output: log.DefaultOutputConfig,
		GlobalVariables: map[string]interface{}{
			"var1": "local",
			"var2": "repo",
			"var3": "local",
			"var4": "global",
		},
	}, actual)
}