package vcs_test

import (
	"fisherman/internal/utils"
	"fisherman/pkg/vcs"
	"fisherman/testing/testutils"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/stretchr/testify/assert"
)

func TestGitRepository_GetFilesInIndex_Empty(t *testing.T) {
	repo, _, fs, w := createRepo()

	testutils.MakeCommits(w, fs, map[string]map[string]string{
		"init commit": {"LICENSE": "MIT"},
		"test commit": {"demo": "this is test file"},
	})

	files, err := repo.GetFilesInIndex()
	assert.NoError(t, err)
	assert.Empty(t, files)
}

func TestGitRepository_GetFilesInIndex_UntrackedFiles(t *testing.T) {
	repo, _, fs, w := createRepo()

	testutils.MakeCommits(w, fs, map[string]map[string]string{
		"init commit": {"LICENSE": "MIT"},
		"test commit": {"demo": "this is test file"},
	})

	testutils.MakeFiles(fs, map[string]string{
		"untracked": "untracked content",
	})

	files, err := repo.GetFilesInIndex()
	assert.NoError(t, err)
	assert.Empty(t, files)
}

func TestGitRepository_GetFilesInIndex(t *testing.T) {
	repo, _, fs, w := createRepo()

	testutils.MakeCommits(w, fs, map[string]map[string]string{
		"init commit": {"LICENSE": "MIT"},
		"test commit": {"demo": "this is test file"},
	})

	testutils.MakeFiles(fs, map[string]string{
		"tracked": "untracked content",
	})

	err := w.AddGlob(".")
	utils.HandleCriticalError(err)

	files, err := repo.GetFilesInIndex()
	assert.NoError(t, err)
	assert.Equal(t, []string{"tracked"}, files)
}

func createRepo() (*vcs.GitRepository, *git.Repository, billy.Filesystem, *git.Worktree) {
	fs := memfs.New()
	r, err := git.Init(memory.NewStorage(), fs)
	utils.HandleCriticalError(err)
	repo := vcs.CreateGitRepository(r)
	utils.HandleCriticalError(err)
	w, err := r.Worktree()
	utils.HandleCriticalError(err)

	return repo, r, fs, w
}
