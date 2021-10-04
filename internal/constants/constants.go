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
	GlobalConfigPath = "GlobalConfigPath"
	LocalConfigPath  = "LocalConfigPath"
	RepoConfigPath   = "RepoConfigPath"
	HookName         = "HookName"
)
