package hooks

func (config *ShellScriptsConfig) GetActive() ScriptsConfig {
	return config.Windows
}

func (config *ShellScriptsConfig) Compile(variables map[string]interface{}) {
	compileCommands(config.Windows, variables)
}
