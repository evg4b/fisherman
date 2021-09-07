package vcs_test

import (
	"fisherman/pkg/guards"
	"fisherman/pkg/vcs"
	"fisherman/testing/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitRepository_GetIndexChanges(t *testing.T) {
	repo, _, fs, w := testutils.CreateRepo(t)

	testutils.MakeCommits(t, w, fs, map[string]map[string]string{
		"init commit":    {"LICENSE": "MIT"},
		"test commit":    {"demo": "this is test file"},
		"test commit 2 ": {"tracked": "this is test file"},
	})

	testutils.MakeFiles(t, fs, map[string]string{
		"tracked":   "this is test file\nadded new content",
		"new":       "new file",
		"untracked": "untracked content",
	})

	err := w.AddGlob("tracked")
	guards.NoError(err)

	err = w.AddGlob("new")
	guards.NoError(err)

	changes, err := repo.GetIndexChanges()

	assert.NoError(t, err)
	assert.Equal(t, map[string]vcs.Changes{
		"new": {
			{Status: vcs.Added, Change: "new file"},
		},
		"tracked": {
			{Status: vcs.Deleted, Change: "this is test file"},
			{Status: vcs.Added, Change: "this is test file\nadded new content"},
		},
	}, changes)
}
