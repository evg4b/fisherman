package hooks

type ScriptConfig struct {
	Commands []string          `yaml:"commands,omitempty"`
	Env      map[string]string `yaml:"env,omitempty"`
	Path     []string          `yaml:"path,omitempty"`
	Output   bool              `yaml:"output,omitempty"`
}

type ScriptsConfig map[string]ScriptConfig

type ShellScriptsConfig struct {
	Windows ScriptsConfig `yaml:"windows,omitempty"`
	Linux   ScriptsConfig `yaml:"linux,omitempty"`
	Darwin  ScriptsConfig `yaml:"darwin,omitempty"`
}

// UnmarshalYAML implements yaml.Unmarshaler interface
func (config *ShellScriptsConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain ShellScriptsConfig
	if err := unmarshal((*plain)(config)); err == nil {
		return nil
	}

	var cmdConfig ScriptsConfig
	err := unmarshal(&cmdConfig)
	if err == nil {
		config.Darwin = cmdConfig
		config.Linux = cmdConfig
		config.Windows = cmdConfig

		return nil
	}

	return err
}

// UnmarshalYAML implements yaml.Unmarshaler interface
func (config *ScriptsConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if *config == nil {
		(*config) = ScriptsConfig{}
	}

	var outputConfig string
	err := unmarshal(&outputConfig)
	if err == nil {
		(*config)["default"] = ScriptConfig{
			Commands: []string{outputConfig},
			Env:      map[string]string{},
		}

		return nil
	}

	var scriptConfig ScriptConfig
	err = unmarshal(&scriptConfig)
	if err == nil {
		(*config)["default"] = scriptConfig

		return nil
	}

	type plain ScriptsConfig
	if err := unmarshal((*plain)(config)); err == nil {
		return nil
	}

	return err
}

// UnmarshalYAML implements yaml.Unmarshaler interface
func (config *ScriptConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var outputConfig string
	err := unmarshal(&outputConfig)
	if err == nil {
		config.Commands = []string{outputConfig}
		config.Env = map[string]string{}

		return nil
	}

	var commands []string
	if err := unmarshal(&commands); err == nil {
		config.Commands = commands
		config.Env = map[string]string{}

		return nil
	}

	type plain ScriptConfig
	if err = unmarshal((*plain)(config)); err == nil {
		return nil
	}

	return err
}
