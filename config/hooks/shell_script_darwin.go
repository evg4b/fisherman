package hooks

func (config *ShellScriptsConfig) GetActive() ScriptsConfig {
	return config.Darwin
}

func (config *ShellScriptsConfig) Compile(variables map[string]interface{}) {
	compileCommands(config.Darwin, variables)
}
