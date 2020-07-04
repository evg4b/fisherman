package git

import (
	"github.com/go-git/go-git/v5"
)

type UserInfo struct {
	Email    string
	UserName string
}

type Snapshot struct {
	repo *git.Repository
}

func NewSnapshot(repo *git.Repository) *Snapshot {
	return &Snapshot{repo}
}
