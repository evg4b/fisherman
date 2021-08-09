package vcs_test

import (
	"fisherman/internal/utils"
	"fisherman/testing/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitRepository_GetFilesInIndex_Empty(t *testing.T) {
	repo, _, fs, w := testutils.CreateRepo()

	testutils.MakeCommits(w, fs, map[string]map[string]string{
		"init commit": {"LICENSE": "MIT"},
		"test commit": {"demo": "this is test file"},
	})

	files, err := repo.GetFilesInIndex()

	assert.NoError(t, err)
	assert.Empty(t, files)
}

func TestGitRepository_GetFilesInIndex_UntrackedFiles(t *testing.T) {
	repo, _, fs, w := testutils.CreateRepo()

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
	repo, _, fs, w := testutils.CreateRepo()

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
