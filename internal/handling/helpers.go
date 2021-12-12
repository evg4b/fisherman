package handling

import (
	"fisherman/internal/configuration"
	"fisherman/internal/constants"
)

func getConfig(name string, config *configuration.HooksConfig) *configuration.HookConfig {
	switch name {
	case constants.ApplyPatchMsgHook:
		return config.ApplyPatchMsgHook
	case constants.CommitMsgHook:
		return config.CommitMsgHook
	case constants.FsMonitorWatchmanHook:
		return config.FsMonitorWatchmanHook
	case constants.PostUpdateHook:
		return config.PostUpdateHook
	case constants.PreApplyPatchHook:
		return config.PreApplyPatchHook
	case constants.PreCommitHook:
		return config.PreCommitHook
	case constants.PrePushHook:
		return config.PrePushHook
	case constants.PreRebaseHook:
		return config.PreRebaseHook
	case constants.PreReceiveHook:
		return config.PreReceiveHook
	case constants.PrepareCommitMsgHook:
		return config.PrepareCommitMsgHook
	case constants.UpdateHook:
		return config.UpdateHook
	}

	panic("incorrect hook name")
}
