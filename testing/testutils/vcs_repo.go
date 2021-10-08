package testutils

import (
	"fisherman/pkg/guards"
	"fisherman/pkg/vcs"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
)

// CreateRepo creates vcs.GitRepository and all dependencies.
func CreateRepo(t *testing.T) (*vcs.GitRepository, *git.Repository, billy.Filesystem, *git.Worktree) {
	t.Helper()

	fs := memfs.New()

	r, err := git.Init(memory.NewStorage(), fs)
	guards.NoError(err)

	repo := vcs.CreateGitRepository(r)
	guards.NoError(err)

	w, err := r.Worktree()
	guards.NoError(err)

	return repo, r, fs, w
}
