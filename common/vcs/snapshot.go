package vcs

type BranchInfo struct {
	Local  string
	Remote string
}

type CommitInfo struct {
	Commit string
}

type Snapshot interface {
	GetHookData() map[string]string
	GetBranchIfo() *BranchInfo
}
