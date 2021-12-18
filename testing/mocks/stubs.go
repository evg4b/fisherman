package mocks

import (
	"fisherman/internal/configuration"
)

var Cwd = "~/repository/project"
var Executable = "string"
var Configs = map[string]string{
	configuration.LocalMode:  "N/A",
	configuration.GlobalMode: "N/A",
	configuration.RepoMode:   "~/repository/project/.fisherman.yaml",
}
var HooksConfigStub = configuration.HooksConfig{}
