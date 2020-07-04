package config

import (
	. "fisherman/config/hooks"
)

type HooksConfig struct {
	ApplyPatchMsgHook     *ApplyPatchMsgHookConfig     `yaml:"applypatch-msg,omitempty"`
	CommitMsgHook         *CommitMsgHookConfig         `yaml:"commit-msg,omitempty"`
	FsMonitorWatchmanHook *FsMonitorWatchmanHookConfig `yaml:"fsmonitor-watchman,omitempty"`
	PostUpdateHook        *PostUpdateHookConfig        `yaml:"post-update,omitempty"`
	PreApplyPatchHook     *PreApplyPatchHookConfig     `yaml:"pre-applypatch,omitempty"`
	PreCommitHook         *PreCommitHookConfig         `yaml:"pre-commit,omitempty"`
	PrePushHook           *PrePushHookConfig           `yaml:"pre-push,omitempty"`
	PreRebaseHook         *PreRebaseHookConfig         `yaml:"pre-rebase,omitempty"`
	PreReceiveHook        *PreReceiveHookConfig        `yaml:"pre-receive,omitempty"`
	PrepareCommitMsgHook  *PrepareCommitMsgHookConfig  `yaml:"prepare-commit-msg,omitempty"`
	UpdateHook            *UpdateHookConfig            `yaml:"update,omitempty"`
}
