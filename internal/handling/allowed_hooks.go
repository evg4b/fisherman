package handling

import (
	"fisherman/constants"
	"fisherman/internal/rules"
)

var allowedHooks = map[string][]string{
	constants.ApplyPatchMsgHook: {
		rules.ShellScriptType,
	},
	constants.CommitMsgHook: {
		rules.ShellScriptType,
		rules.CommitMessageType,
	},
	constants.FsMonitorWatchmanHook: {
		rules.ShellScriptType,
	},
	constants.PostUpdateHook: {
		rules.ShellScriptType,
	},
	constants.PreApplyPatchHook: {
		rules.ShellScriptType,
	},
	constants.PreCommitHook: {
		rules.ShellScriptType,
		rules.AddToIndexType,
		rules.SuppressCommitFilesType,
	},
	constants.PrePushHook: {
		rules.ShellScriptType,
	},
	constants.PreRebaseHook: {
		rules.ShellScriptType,
	},
	constants.PreReceiveHook: {
		rules.ShellScriptType,
	},
	constants.PrepareCommitMsgHook: {
		rules.PrepareMessageType,
	},
	constants.UpdateHook: {
		rules.ShellScriptType,
	},
}
