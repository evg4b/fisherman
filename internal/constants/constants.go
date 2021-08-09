package constants

const AppName = "fisherman"

var Version = "x.x.x"

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
