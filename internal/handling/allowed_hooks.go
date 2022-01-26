package handling

import (
	"fisherman/internal/constants"
	"fisherman/internal/rules"
)

var allowedHooks = map[string][]string{
	constants.ApplyPatchMsgHook: {
		rules.ShellScriptType,
		rules.ExecType,
	},
	constants.CommitMsgHook: {
		rules.ShellScriptType,
		rules.ExecType,
		rules.CommitMessageType,
	},
	constants.FsMonitorWatchmanHook: {
		rules.ShellScriptType,
		rules.ExecType,
	},
	constants.PreApplyPatchHook: {
		rules.ShellScriptType,
		rules.ExecType,
	},
	constants.PreCommitHook: {
		rules.ShellScriptType,
		rules.ExecType,
		rules.AddToIndexType,
		rules.SuppressCommitFilesType,
		rules.SuppressedTextType,
	},
	constants.PrePushHook: {
		rules.ShellScriptType,
		rules.ExecType,
	},
	constants.PreRebaseHook: {
		rules.ShellScriptType,
		rules.ExecType,
	},
	constants.PrepareCommitMsgHook: {
		rules.PrepareMessageType,
	},
}
