package mocks

import (
	"fisherman/configuration"
	"fisherman/internal"
)

var AppInfoStub = internal.AppInfo{
	Cwd:        "~/repository/project",
	Executable: "string",
	Configs: map[string]string{
		configuration.LocalMode:  "N/A",
		configuration.GlobalMode: "N/A",
		configuration.RepoMode:   "~/repository/project/.fisherman.yaml",
	},
}

var HooksConfigStub = configuration.HooksConfig{}
