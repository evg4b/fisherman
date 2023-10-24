package mocks

import (
	"github.com/evg4b/fisherman/internal/configuration"
)

var (
	Cwd        = "~/repository/project"
	Executable = "string"
	Configs    = map[string]string{
		configuration.LocalMode:  "N/A",
		configuration.GlobalMode: "N/A",
		configuration.RepoMode:   "~/repository/project/.fisherman.yaml",
	}
)
var HooksConfigStub = configuration.HooksConfig{}
