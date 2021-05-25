package mocks

import (
	"fisherman/internal"
	"fisherman/internal/configuration"
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
