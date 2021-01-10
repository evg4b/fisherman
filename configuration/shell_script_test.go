package configuration_test

import (
	. "fisherman/configuration" // nolint
	"fisherman/constants"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestShellScriptFullStructure_UnmarshalYAML(t *testing.T) {
	var yamlMarkup string = `
windows:
  script1:
    env:
      demo: Windows
    commands:
      - windowsCommand1
      - windowsCommand2
linux:
  script1:
    env:
      demo: Linux
    commands:
      - linuxCommand1
      - linuxCommand2
darwin:
  script1:
    env:
      demo: Darwin
    commands:
      - darwinCommand1
      - darwinCommand2
`
	var data ScriptsConfig
	err := decode(yamlMarkup, &data)

	expectedLinuxConfig := ScriptsConfig{
		"script1": ScriptConfig{
			Commands: []string{"linuxCommand1", "linuxCommand2"},
			Env: map[string]string{
				"demo": "Linux",
			},
		},
	}

	expectedWindowsConfig := ScriptsConfig{
		"script1": ScriptConfig{
			Commands: []string{"windowsCommand1", "windowsCommand2"},
			Env: map[string]string{
				"demo": "Windows",
			},
		},
	}

	expectedDarwinConfig := ScriptsConfig{
		"script1": ScriptConfig{
			Commands: []string{"darwinCommand1", "darwinCommand2"},
			Env: map[string]string{
				"demo": "Darwin",
			},
		},
	}

	assert.NoError(t, err)
	switch runtime.GOOS {
	case constants.LinuxOS:
		assert.Equal(t, expectedLinuxConfig, data)
	case constants.WindowsOS:
		assert.Equal(t, expectedWindowsConfig, data)
	case constants.DarwinOS:
		assert.Equal(t, expectedDarwinConfig, data)
	}
}

func TestShellScriptWithoutSystem_UnmarshalYAML(t *testing.T) {
	var yamlMarkup string = `
script1:
  output: true
  env:
    demo: demo1
  commands:
    - command1
script2:
  env:
    demo: demo2
  commands:
    - command1
    - command2
`
	var data ScriptsConfig
	err := decode(yamlMarkup, &data)

	expectedConfig := ScriptsConfig{
		"script1": ScriptConfig{
			Output:   true,
			Commands: []string{"command1"},
			Env:      map[string]string{"demo": "demo1"},
		},
		"script2": ScriptConfig{
			Commands: []string{"command1", "command2"},
			Env:      map[string]string{"demo": "demo2"},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedConfig, data)
}

func TestShellScriptShortString_UnmarshalYAML(t *testing.T) {
	var yamlMarkup string = `
test1: command1
test2: command2
test3: command3
`
	var data ScriptsConfig
	err := decode(yamlMarkup, &data)

	expectedConfig := ScriptsConfig{
		"test1": ScriptConfig{Commands: []string{"command1"}, Env: map[string]string{}},
		"test2": ScriptConfig{Commands: []string{"command2"}, Env: map[string]string{}},
		"test3": ScriptConfig{Commands: []string{"command3"}, Env: map[string]string{}},
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedConfig, data)
}

func TestShellScriptOneScriptOnly_UnmarshalYAML(t *testing.T) {
	var yamlMarkup string = `
env:
  demo: demo2
commands:
  - command1
  - command2
`
	var data ScriptsConfig
	err := decode(yamlMarkup, &data)

	expectedConfig := ScriptsConfig{
		"default": ScriptConfig{
			Commands: []string{"command1", "command2"},
			Env:      map[string]string{"demo": "demo2"},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedConfig, data)
}

func TestShellScriptOneScrtipOnly_UnmarshalYAML(t *testing.T) {
	var yamlMarkup string = `demo-command`
	var data ScriptsConfig
	err := decode(yamlMarkup, &data)

	expectedConfig := ScriptsConfig{
		"default": ScriptConfig{
			Commands: []string{"demo-command"},
			Env:      map[string]string{},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedConfig, data)
}

func TestShellScriptStringArray_UnmarshalYAML(t *testing.T) {
	var yamlMarkup string = `
- demo-command1
- demo-command2
- demo-command3
`
	var data ScriptsConfig
	err := decode(yamlMarkup, &data)

	expectedConfig := ScriptsConfig{
		"default": ScriptConfig{
			Commands: []string{"demo-command1", "demo-command2", "demo-command3"},
			Env:      map[string]string{},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedConfig, data)
}

func TestShellScriptStringArrayWithBrackets_UnmarshalYAML(t *testing.T) {
	var yamlMarkup string = `[cmd1, cmd2, cmd3]`
	var data ScriptsConfig
	err := decode(yamlMarkup, &data)

	expectedConfig := ScriptsConfig{
		"default": ScriptConfig{
			Commands: []string{"cmd1", "cmd2", "cmd3"},
			Env:      map[string]string{},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedConfig, data)
}

func decode(yamlMarkup string, v interface{}) error {
	decoder := yaml.NewDecoder(strings.NewReader(yamlMarkup))
	decoder.KnownFields(true)

	return decoder.Decode(v)
}
