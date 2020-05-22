package configuration

type HooksConfig struct {
	ApplypatchMsgHook        *HookConfig `yaml:"applypatch-msg"`
	PreApplypatchHook        *HookConfig `yaml:"pre-applypatch"`
	PostApplypatchHook       *HookConfig `yaml:"post-applypatch"`
	PreCommitHook            *HookConfig `yaml:"pre-commit"`
	PreMergeCommitHook       *HookConfig `yaml:"pre-merge-commit"`
	PrepareCommitMsgHook     *HookConfig `yaml:"prepare-commit-msg"`
	CommitMsgHook            *HookConfig `yaml:"commit-msg"`
	PostCommitHook           *HookConfig `yaml:"post-commit"`
	PreRebaseHook            *HookConfig `yaml:"pre-rebase"`
	PostCheckoutHook         *HookConfig `yaml:"post-checkout"`
	PostMergeHook            *HookConfig `yaml:"post-merge"`
	PrePushHook              *HookConfig `yaml:"pre-push"`
	PreReceiveHook           *HookConfig `yaml:"pre-receive"`
	UpdateHook               *HookConfig `yaml:"update"`
	ProcReceiveHook          *HookConfig `yaml:"proc-receive"`
	PostReceiveHook          *HookConfig `yaml:"post-receive"`
	PostUpdateHook           *HookConfig `yaml:"post-update"`
	ReferenceTransactionHook *HookConfig `yaml:"reference-transaction"`
	PushToCheckoutHook       *HookConfig `yaml:"push-to-checkout"`
	PreAutoGcHook            *HookConfig `yaml:"pre-auto-gc"`
	PostRewriteHook          *HookConfig `yaml:"post-rewrite"`
	SendemailValidateHook    *HookConfig `yaml:"sendemail-validate"`
	FsmonitorWatchmanHook    *HookConfig `yaml:"fsmonitor-watchman"`
	P4ChangelistHook         *HookConfig `yaml:"p4-changelist"`
	P4PrepareChangelistHook  *HookConfig `yaml:"p4-prepare-changelist"`
	P4PostChangelistHook     *HookConfig `yaml:"p4-post-changelist"`
	P4PreSubmitHook          *HookConfig `yaml:"p4-pre-submit"`
	PostIndexChangeHook      *HookConfig `yaml:"post-index-change"`
}
