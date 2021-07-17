package vcs

import (
	"fisherman/testing/testutils"
	"testing"

	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/stretchr/testify/assert"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
)

func TestGitRepository_GetFilesInIndex_Empty(t *testing.T) {
	fs := memfs.New()
	r, _ := git.Init(memory.NewStorage(), fs)

	w, _ := r.Worktree()
	testutils.MakeCommits(w, fs, map[string]map[string]string{
		"init commit": {"LICENSE": "MIT"},
		"test commit": {"demo": "this is test file"},
	})

	repo := OpenGitRepository("test")
	_, _ = repo.repo()
	repo.internalRepo = r

	files, err := repo.GetFilesInIndex()
	assert.NoError(t, err)
	assert.Empty(t, files)
}

func TestGitRepository_GetFilesInIndex_UntrackedFiles(t *testing.T) {
	fs := memfs.New()
	r, _ := git.Init(memory.NewStorage(), fs)

	w, _ := r.Worktree()
	testutils.MakeCommits(w, fs, map[string]map[string]string{
		"init commit": {"LICENSE": "MIT"},
		"test commit": {"demo": "this is test file"},
	})

	testutils.MakeFiles(fs, map[string]string{
		"untracked": "untracked content",
	})

	repo := OpenGitRepository("test")
	_, _ = repo.repo()
	repo.internalRepo = r

	files, err := repo.GetFilesInIndex()
	assert.NoError(t, err)
	assert.Empty(t, files)
}

func TestGitRepository_GetFilesInIndex(t *testing.T) {
	fs := memfs.New()
	r, _ := git.Init(memory.NewStorage(), fs)

	w, _ := r.Worktree()
	testutils.MakeCommits(w, fs, map[string]map[string]string{
		"init commit": {"LICENSE": "MIT"},
		"test commit": {"demo": "this is test file"},
	})

	testutils.MakeFiles(fs, map[string]string{
		"tracked": "untracked content",
	})

	_ = w.AddGlob(".")

	repo := OpenGitRepository("test")
	_, _ = repo.repo()
	repo.internalRepo = r

	files, err := repo.GetFilesInIndex()
	assert.NoError(t, err)
	assert.Equal(t, []string{"tracked"}, files)
}
