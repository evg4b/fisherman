package expression

import (
	"github.com/evg4b/fisherman/internal/constants"
	"github.com/evg4b/fisherman/internal/utils"
)

type EnvVars map[string]any

func (vars EnvVars) IsEmpty(value string) bool {
	return utils.IsEmpty(value)
}

func (vars EnvVars) IsWindows() bool {
	return vars[constants.OsVariable] == Windows
}

func (vars EnvVars) IsLinux() bool {
	return vars[constants.OsVariable] == Linux
}

func (vars EnvVars) IsMacOs() bool {
	return vars[constants.OsVariable] == Macos
}
