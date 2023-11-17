package handling

import (
	"github.com/evg4b/fisherman/internal/constants"
	"github.com/evg4b/fisherman/internal/rules"
)

var allowedHooks = map[string][]string{
	constants.PostApplypatchHook: {
		rules.ShellScriptType,
		rules.ExecType,
	},
	constants.PreMergeCommitHook: {
		rules.ShellScriptType,
		rules.ExecType,
	},
	constants.PostCommitHook: {
		rules.ShellScriptType,
		rules.ExecType,
	},
	constants.PostCheckoutHook: {
		rules.ShellScriptType,
		rules.ExecType,
	},
	constants.PostMergeHook: {
		rules.ShellScriptType,
		rules.ExecType,
	},
	constants.PreReceiveHook: {
		rules.ShellScriptType,
		rules.ExecType,
	},
	constants.UpdateHook: {
		rules.ShellScriptType,
		rules.ExecType,
	},
	constants.ProcReceiveHook: {
		rules.ShellScriptType,
		rules.ExecType,
	},
	constants.PostReceiveHook: {
		rules.ShellScriptType,
		rules.ExecType,
	},
	constants.PostUpdateHook: {
		rules.ShellScriptType,
		rules.ExecType,
	},
	constants.ReferenceTransactionHook: {
		rules.ShellScriptType,
		rules.ExecType,
	},
	constants.PushToCheckoutHook: {
		rules.ShellScriptType,
		rules.ExecType,
	},
	constants.PreAutoGcHook: {
		rules.ShellScriptType,
		rules.ExecType,
	},
	constants.PostRewriteHook: {
		rules.ShellScriptType,
		rules.ExecType,
	},
	constants.SendemailValidateHook: {
		rules.ShellScriptType,
		rules.ExecType,
	},
	constants.FsmonitorWatchmanHook: {
		rules.ShellScriptType,
		rules.ExecType,
	},
	constants.P4ChangelistHook: {
		rules.ShellScriptType,
		rules.ExecType,
	},
	constants.P4PrepareChangelistHook: {
		rules.ShellScriptType,
		rules.ExecType,
	},
	constants.P4PostChangelistHook: {
		rules.ShellScriptType,
		rules.ExecType,
	},
	constants.P4PreSubmitHook: {
		rules.ShellScriptType,
		rules.ExecType,
	},
	constants.PostIndexChangeHook: {
		rules.ShellScriptType,
		rules.ExecType,
	},
	constants.ApplypatchMsgHook: {
		rules.ShellScriptType,
		rules.ExecType,
	},
	constants.CommitMsgHook: {
		rules.ShellScriptType,
		rules.ExecType,
		rules.CommitMessageType,
	},
	constants.FsmonitorWatchmanHook: {
		rules.ShellScriptType,
		rules.ExecType,
	},
	constants.PreApplypatchHook: {
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
