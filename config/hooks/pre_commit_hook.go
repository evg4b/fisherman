package hooks

import "fisherman/utils"

// PreCommitHookConfig is structure to storage user configuration for pre-commit hook
type PreCommitHookConfig struct {
	Cmd ScriptsConfig
}

func (config *PreCommitHookConfig) Compile(variables map[string]interface{}) {
	for _, cmd := range config.Cmd {
		for key := range cmd.Commands {
			utils.FillTemplate(&cmd.Commands[key], variables)
		}
	}
}
