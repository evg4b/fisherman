package handling

import (
	"fisherman/internal/constants"
	"fisherman/internal/rules"
)

var allowedHooks = map[string][]string{
	constants.ApplyPatchMsgHook: {
		rules.ShellScriptType,
		rules.RunProgramType,
	},
	constants.CommitMsgHook: {
		rules.ShellScriptType,
		rules.RunProgramType,
		rules.CommitMessageType,
	},
	constants.FsMonitorWatchmanHook: {
		rules.ShellScriptType,
		rules.RunProgramType,
	},
	constants.PostUpdateHook: {
		rules.ShellScriptType,
		rules.RunProgramType,
	},
	constants.PreApplyPatchHook: {
		rules.ShellScriptType,
		rules.RunProgramType,
	},
	constants.PreCommitHook: {
		rules.ShellScriptType,
		rules.RunProgramType,
		rules.AddToIndexType,
		rules.SuppressCommitFilesType,
	},
	constants.PrePushHook: {
		rules.ShellScriptType,
		rules.RunProgramType,
	},
	constants.PreRebaseHook: {
		rules.ShellScriptType,
		rules.RunProgramType,
	},
	constants.PreReceiveHook: {
		rules.ShellScriptType,
		rules.RunProgramType,
	},
	constants.PrepareCommitMsgHook: {
		rules.PrepareMessageType,
	},
	constants.UpdateHook: {
		rules.ShellScriptType,
		rules.RunProgramType,
	},
}
