package config

var (
	// GlobalMode is constant for storage config in user folder `~/.fisherman.yaml`
	GlobalMode = "global"
	// LocalMode is constant for storage config in git folder `<repo>/.git/hooks/.fisherman.yaml`
	LocalMode = "local"
	// RepoMode is constant for storage config in repo folder `<repo>/.fisherman.yaml`
	RepoMode = "repo"
)
