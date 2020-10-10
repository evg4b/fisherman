package hooks

type Command struct {
	Commands []string
	Env      map[string]string
}

type CmdConfig map[string]Command

// UnmarshalYAML implements yaml.Unmarshaler interface
func (config *CmdConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var outputConfig string
	err := unmarshal(&outputConfig)
	if err == nil {
		(*config)["cmd"] = Command{
			Commands: []string{outputConfig},
			Env:      map[string]string{},
		}

		return nil
	}

	type plain CmdConfig
	if err := unmarshal((*plain)(config)); err != nil {
		return nil
	}

	return nil
}

// UnmarshalYAML implements yaml.Unmarshaler interface
func (config *Command) UnmarshalYAML(unmarshal func(interface{}) error) error {
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

	type plain Command
	if err = unmarshal((*plain)(config)); err != nil {
		return err
	}

	return nil
}
