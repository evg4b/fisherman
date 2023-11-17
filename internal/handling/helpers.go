package handling

import (
	"fmt"
	"github.com/evg4b/fisherman/internal/configuration"
	"github.com/evg4b/fisherman/internal/constants"
)

// nolint: cyclop
func getConfig(name string, config *configuration.HooksConfig) (*configuration.HookConfig, error) {
	switch name {
	case constants.ApplypatchMsgHook:
		return config.ApplypatchMsgHook, nil
	case constants.PreApplypatchHook:
		return config.PreApplypatchHook, nil
	case constants.PostApplypatchHook:
		return config.PostApplypatchHook, nil
	case constants.PreCommitHook:
		return config.PreCommitHook, nil
	case constants.PreMergeCommitHook:
		return config.PreMergeCommitHook, nil
	case constants.PrepareCommitMsgHook:
		return config.PrepareCommitMsgHook, nil
	case constants.CommitMsgHook:
		return config.CommitMsgHook, nil
	case constants.PostCommitHook:
		return config.PostCommitHook, nil
	case constants.PreRebaseHook:
		return config.PreRebaseHook, nil
	case constants.PostCheckoutHook:
		return config.PostCheckoutHook, nil
	case constants.PostMergeHook:
		return config.PostMergeHook, nil
	case constants.PrePushHook:
		return config.PrePushHook, nil
	case constants.PreReceiveHook:
		return config.PreReceiveHook, nil
	case constants.UpdateHook:
		return config.UpdateHook, nil
	case constants.ProcReceiveHook:
		return config.ProcReceiveHook, nil
	case constants.PostReceiveHook:
		return config.PostReceiveHook, nil
	case constants.PostUpdateHook:
		return config.PostUpdateHook, nil
	case constants.ReferenceTransactionHook:
		return config.ReferenceTransactionHook, nil
	case constants.PushToCheckoutHook:
		return config.PushToCheckoutHook, nil
	case constants.PreAutoGcHook:
		return config.PreAutoGcHook, nil
	case constants.PostRewriteHook:
		return config.PostRewriteHook, nil
	case constants.SendemailValidateHook:
		return config.SendemailValidateHook, nil
	case constants.FsmonitorWatchmanHook:
		return config.FsmonitorWatchmanHook, nil
	case constants.P4ChangelistHook:
		return config.P4ChangelistHook, nil
	case constants.P4PrepareChangelistHook:
		return config.P4PrepareChangelistHook, nil
	case constants.P4PostChangelistHook:
		return config.P4PostChangelistHook, nil
	case constants.P4PreSubmitHook:
		return config.P4PreSubmitHook, nil
	case constants.PostIndexChangeHook:
		return config.PostIndexChangeHook, nil
	}

	return nil, fmt.Errorf("'%s' is not valid hook name", name)
}
