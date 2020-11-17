package hooks

type AddToIndexConfig struct {
	Globs    []string `yaml:"globs,omitempty"`
	Optional bool     `yaml:"optional,omitempty"`
}

func (config *AddToIndexConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var globs []string
	err := unmarshal(&globs)
	if err == nil {
		(*config) = AddToIndexConfig{
			Globs:    globs,
			Optional: false,
		}

		return nil
	}

	type plain AddToIndexConfig

	return unmarshal((*plain)(config))
}

type PreCommitHookConfig struct {
	Variables       Variables        `yaml:"variables,omitempty"`
	Shell           ScriptsConfig    `yaml:"shell,omitempty"`
	AddFilesToIndex AddToIndexConfig `yaml:"add-to-index,omitempty"`
}

func (config *PreCommitHookConfig) Compile(variables map[string]interface{}) {
	config.Shell.Compile(variables)
}

func (config *PreCommitHookConfig) GetVarsSection() Variables {
	return config.Variables
}

func (*PreCommitHookConfig) HasVars() bool {
	return true
}

func (config *PreCommitHookConfig) IsEmpty() bool {
	return len(config.Shell) == 0 &&
		len(config.AddFilesToIndex.Globs) == 0 &&
		!config.AddFilesToIndex.Optional &&
		config.Variables == Variables{}
}
