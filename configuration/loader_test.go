package configuration_test

import (
	"fisherman/configuration"
	"fisherman/constants"
	"fisherman/infrastructure/log"
	"fisherman/internal/rules"
	"fisherman/testing/mocks"
	"fisherman/testing/testutils"
	"fmt"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigLoader_FindConfigFiles(t *testing.T) {
	usr := user.User{HomeDir: filepath.Join("/", "usr", "home")}
	cwd := filepath.Join("/", "usr", "home", "documents", "repo")

	localConfig := configuration.GetConfigFolder(&usr, cwd, configuration.LocalMode)
	repoConfig := configuration.GetConfigFolder(&usr, cwd, configuration.RepoMode)
	globalConfig := configuration.GetConfigFolder(&usr, cwd, configuration.GlobalMode)

	tests := []struct {
		name        string
		files       []string
		expected    map[string]string
		expectedErr error
	}{
		{
			name: "mere then one config file",
			files: []string{
				filepath.Join(localConfig, constants.AppConfigNames[0]),
				filepath.Join(localConfig, constants.AppConfigNames[1]),
				filepath.Join(repoConfig, constants.AppConfigNames[0]),
				filepath.Join(globalConfig, constants.AppConfigNames[0]),
			},
			expectedErr: fmt.Errorf("more then one config file specifies in folder '%s'", localConfig),
		},
		{
			name: "correct files loading",
			files: []string{
				filepath.Join(localConfig, constants.AppConfigNames[0]),
				filepath.Join(repoConfig, constants.AppConfigNames[0]),
				filepath.Join(globalConfig, constants.AppConfigNames[0]),
			},
			expected: map[string]string{
				configuration.LocalMode:  filepath.Join(localConfig, constants.AppConfigNames[0]),
				configuration.RepoMode:   filepath.Join(repoConfig, constants.AppConfigNames[0]),
				configuration.GlobalMode: filepath.Join(globalConfig, constants.AppConfigNames[0]),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fs = testutils.FsFromSlice(t, tt.files)
			loaded := configuration.NewLoader(&usr, cwd, fs)

			actual, err := loaded.FindConfigFiles()

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

	fs := testutils.FsFromMap(t, map[string]string{
		"GlobalConfig":      config,
		"GlobalConfigError": "asd['",
	})

	tests := []struct {
		name        string
		loader      *configuration.ConfigLoader
		files       map[string]string
		expected    *configuration.FishermanConfig
		expectedErr string
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
			name:   "file reader error",
			loader: configuration.NewLoader(&usr, cwd, fs),
			files: map[string]string{
				configuration.GlobalMode: "GlobalConfig3",
			},
			expectedErr: "open GlobalConfig3: file does not exist",
		},
		{
			name:   "correct decoding",
			loader: configuration.NewLoader(&usr, cwd, fs),
			files: map[string]string{
				configuration.GlobalMode: "GlobalConfig",
			},
			expected: &configuration.FishermanConfig{
				Output: log.DefaultOutputConfig,
				GlobalVariables: map[string]interface{}{
					"name": "value",
				},
				Hooks: configuration.HooksConfig{
					PrePushHook: &configuration.HookConfig{
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
			name:   "decoding error",
			loader: configuration.NewLoader(&usr, cwd, fs),
			files: map[string]string{
				configuration.GlobalMode: "GlobalConfigError",
			},
			expectedErr: "yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `asd['` into configuration.FishermanConfig",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := tt.loader.Load(tt.files)

			testutils.CheckError(t, tt.expectedErr, err)
			if tt.expectedErr == "" {
				assert.Equal(t, tt.expected, actual)
			}
		})
	}
}

func TestConfigLoader_Load_Correct_Merging(t *testing.T) {
	usr := user.User{HomeDir: filepath.Join("/", "usr", "home")}
	cwd := filepath.Join("/", "usr", "home", "documents", "repo")

	fs := testutils.FsFromMap(t, map[string]string{
		"global.yaml": `
variables:
  var1: global
  var2: global
  var3: global
  var4: global`,
		"repo.yaml": `
variables:
  var1: repo
  var2: repo`,
		"local.yaml": `
variables:
  var1: local
  var3: local`,
	})

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
