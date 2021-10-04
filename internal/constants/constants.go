package constants

// AppName is a binary file name.
const AppName = "fisherman"

// Version is version of binary.
var Version = "x.x.x"

// AppConfigNames is allowed config names.
var AppConfigNames = []string{
	".fisherman.yaml",
	".fisherman.yml",
}

const (
	// GlobalConfigPath is the identifier for the global path of the config file.
	// This file is located in user folder '~/fisherman.yaml'.
	GlobalConfigPath = "GlobalConfigPath"
	// LocalConfigPath is the identifier for the local path of the config file.
	// This configuration file is located in the .git folder of the repository.
	// This file will not be included in the repository index.
	LocalConfigPath = "LocalConfigPath"
	// RepoConfigPath is the identifier for the repository level path of the config file.
	// This configuration file is located in the repository and it will be included in the repository index.
	RepoConfigPath = "RepoConfigPath"
	// HookName is constant for identification name of current hook in template.
	HookName = "HookName"
)
