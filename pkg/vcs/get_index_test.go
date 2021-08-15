package vcs_test

import (
	"fisherman/internal/utils"
	"fisherman/pkg/vcs"
	"fisherman/testing/testutils"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
)

func TestGitRepository_GetFilesInIndex_Empty(t *testing.T) {
	repo, _, fs, w := testutils.CreateRepo(t)

	testutils.MakeCommits(t, w, fs, map[string]map[string]string{
		"init commit": {"LICENSE": "MIT"},
		"test commit": {"demo": "this is test file"},
	})

	files, err := repo.GetFilesInIndex()

	assert.NoError(t, err)
	assert.Empty(t, files)
}

func TestGitRepository_GetFilesInIndex_UntrackedFiles(t *testing.T) {
	repo, _, fs, w := testutils.CreateRepo(t)

	testutils.MakeCommits(t, w, fs, map[string]map[string]string{
		"init commit": {"LICENSE": "MIT"},
		"test commit": {"demo": "this is test file"},
	})

	testutils.MakeFiles(t, fs, map[string]string{
		"untracked": "untracked content",
	})

	files, err := repo.GetFilesInIndex()

	assert.NoError(t, err)
	assert.Empty(t, files)
}

func TestGitRepository_GetFilesInIndex(t *testing.T) {
	repo, _, fs, w := testutils.CreateRepo(t)

	testutils.MakeCommits(t, w, fs, map[string]map[string]string{
		"init commit": {"LICENSE": "MIT"},
		"test commit": {"demo": "this is test file"},
	})

	testutils.MakeFiles(t, fs, map[string]string{
		"tracked": "tracked content",
	})

	err := w.AddGlob(".")
	utils.HandleCriticalError(err)

	files, err := repo.GetFilesInIndex()

	assert.NoError(t, err)
	assert.Equal(t, []string{"tracked"}, files)
}

func TestGitRepository_GetIndexChanges(t *testing.T) {
	t.SkipNow()

	repo, _, fs, w := testutils.CreateRepo(t)

	testutils.MakeCommits(t, w, fs, map[string]map[string]string{
		"init commit": {"LICENSE": "MIT"},
		"test commit": {"demo": "this is test file"},
	})

	testutils.MakeFiles(t, fs, map[string]string{
		"tracked":   "tracked content",
		"untracked": "untracked content",
	})

	err := w.AddGlob(".")
	utils.HandleCriticalError(err)

	changes, err := repo.GetIndexChanges()

	assert.NoError(t, err)
	assert.Equal(t, map[string]vcs.Changes{
		"tracked": {
			{Status: git.Added, Change: "untracked content"},
		},
	}, changes)
}
