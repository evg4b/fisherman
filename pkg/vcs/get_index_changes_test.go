// nolint: dupl
package vcs_test

import (
	"github.com/evg4b/fisherman/pkg/guards"
	"github.com/evg4b/fisherman/testing/testutils"
	"testing"

	. "github.com/evg4b/fisherman/pkg/vcs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGitRepository_GetIndexChanges(t *testing.T) {
	t.Run("returns files", func(t *testing.T) {
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

		require.NoError(t, err)
		assert.Equal(t, map[string]Changes{
			"new": {
				{Status: Added, Change: "new file"},
			},
			"tracked": {
				{Status: Deleted, Change: "this is test file"},
				{Status: Added, Change: "this is test file\nadded new content"},
			},
		}, changes)
	})

	t.Run("returns files include subfolders", func(t *testing.T) {
		repo, _, fs, w := testutils.CreateRepo(t)

		testutils.MakeCommits(t, w, fs, map[string]map[string]string{
			"init commit":    {"folder1/config.json": "{}"},
			"test commit":    {"folder1/demo": "this is test file"},
			"test commit 2 ": {"folder1/existing": "this is test file"},
		})

		testutils.MakeFiles(t, fs, map[string]string{
			"folder1/existing": "this is test file\nadded new content",
			"folder1/added":    "added content",
			"untracked":        "untracked content",
		})

		err := w.AddGlob("folder1/existing")
		guards.NoError(err)

		err = w.AddGlob("folder1/added")
		guards.NoError(err)

		changes, err := repo.GetIndexChanges()

		require.NoError(t, err)
		assert.Equal(t, map[string]Changes{
			"folder1/added": {
				{Status: Added, Change: "added content"},
			},
			"folder1/existing": {
				{Status: Deleted, Change: "this is test file"},
				{Status: Added, Change: "this is test file\nadded new content"},
			},
		}, changes)
	})

	t.Run("no changes", func(t *testing.T) {
		repo, _, fs, w := testutils.CreateRepo(t)

		testutils.MakeCommits(t, w, fs, map[string]map[string]string{
			"init commit": {"LICENSE": "MIT"},
		})

		changes, err := repo.GetIndexChanges()

		require.NoError(t, err)
		assert.Equal(t, map[string]Changes{}, changes)
	})

	t.Run("untracked changes", func(t *testing.T) {
		repo, _, fs, w := testutils.CreateRepo(t)

		testutils.MakeCommits(t, w, fs, map[string]map[string]string{
			"init commit": {"LICENSE": "MIT"},
		})

		testutils.MakeFiles(t, fs, map[string]string{
			"file": "content",
		})

		changes, err := repo.GetIndexChanges()

		require.NoError(t, err)
		assert.Equal(t, map[string]Changes{}, changes)
	})
}
