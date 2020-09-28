package config

import (
	"fisherman/config/hooks"
)

// HooksConfig is structure for storage information for hook validation
type HooksConfig struct {
	ApplyPatchMsgHook     *hooks.ApplyPatchMsgHookConfig     `yaml:"applypatch-msg,omitempty"`
	CommitMsgHook         *hooks.CommitMsgHookConfig         `yaml:"commit-msg,omitempty"`
	FsMonitorWatchmanHook *hooks.FsMonitorWatchmanHookConfig `yaml:"fsmonitor-watchman,omitempty"`
	PostUpdateHook        *hooks.PostUpdateHookConfig        `yaml:"post-update,omitempty"`
	PreApplyPatchHook     *hooks.PreApplyPatchHookConfig     `yaml:"pre-applypatch,omitempty"`
	PreCommitHook         *hooks.PreCommitHookConfig         `yaml:"pre-commit,omitempty"`
	PrePushHook           *hooks.PrePushHookConfig           `yaml:"pre-push,omitempty"`
	PreRebaseHook         *hooks.PreRebaseHookConfig         `yaml:"pre-rebase,omitempty"`
	PreReceiveHook        *hooks.PreReceiveHookConfig        `yaml:"pre-receive,omitempty"`
	PrepareCommitMsgHook  *hooks.PrepareCommitMsgHookConfig  `yaml:"prepare-commit-msg,omitempty"`
	UpdateHook            *hooks.UpdateHookConfig            `yaml:"update,omitempty"`
}
