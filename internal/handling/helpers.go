package handling

import (
	"fisherman/internal/configuration"
	"fisherman/internal/constants"
	"fmt"
)

// nolint: cyclop
func getConfig(name string, config *configuration.HooksConfig) (*configuration.HookConfig, error) {
	switch name {
	case constants.ApplyPatchMsgHook:
		return config.ApplyPatchMsgHook, nil
	case constants.CommitMsgHook:
		return config.CommitMsgHook, nil
	case constants.FsMonitorWatchmanHook:
		return config.FsMonitorWatchmanHook, nil
	case constants.PostUpdateHook:
		return config.PostUpdateHook, nil
	case constants.PreApplyPatchHook:
		return config.PreApplyPatchHook, nil
	case constants.PreCommitHook:
		return config.PreCommitHook, nil
	case constants.PrePushHook:
		return config.PrePushHook, nil
	case constants.PreRebaseHook:
		return config.PreRebaseHook, nil
	case constants.PreReceiveHook:
		return config.PreReceiveHook, nil
	case constants.PrepareCommitMsgHook:
		return config.PrepareCommitMsgHook, nil
	case constants.UpdateHook:
		return config.UpdateHook, nil
	}

	return nil, fmt.Errorf("'%s' is not valid hook name", name)
}
