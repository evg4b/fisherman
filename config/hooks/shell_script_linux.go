package hooks

func (config *ShellScriptsConfig) GetActive() ScriptsConfig {
	return config.Linux
}

func (config *ShellScriptsConfig) Compile(variables map[string]interface{}) {
	compileCommands(config.Linux, variables)
}
