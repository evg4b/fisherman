package testutils

import (
	"fisherman/internal/utils"
	"fisherman/pkg/vcs"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
)

func CreateRepo() (*vcs.GitRepository, *git.Repository, billy.Filesystem, *git.Worktree) {
	fs := memfs.New()
	r, err := git.Init(memory.NewStorage(), fs)
	utils.HandleCriticalError(err)

	repo := vcs.CreateGitRepository(r)
	utils.HandleCriticalError(err)

	w, err := r.Worktree()
	utils.HandleCriticalError(err)

	return repo, r, fs, w
}
