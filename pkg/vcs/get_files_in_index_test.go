package vcs_test

import (
	"errors"
	"github.com/evg4b/fisherman/pkg/guards"
	"github.com/evg4b/fisherman/testing/mocks"
	"github.com/evg4b/fisherman/testing/testutils"
	"testing"

	. "github.com/evg4b/fisherman/pkg/vcs"

	"github.com/go-git/go-git/v5/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGitRepository_GetFilesInIndex(t *testing.T) {
	t.Run("no files", func(t *testing.T) {
		repo, _, fs, w := testutils.CreateRepo(t)

		testutils.MakeCommits(t, w, fs, map[string]map[string]string{
			"init commit": {"LICENSE": "MIT"},
			"test commit": {"demo": "this is test file"},
		})

		files, err := repo.GetFilesInIndex()

		require.NoError(t, err)
		assert.Empty(t, files)
	})

	t.Run("excluded untracked files", func(t *testing.T) {
		repo, _, fs, w := testutils.CreateRepo(t)

		testutils.MakeCommits(t, w, fs, map[string]map[string]string{
			"init commit": {"LICENSE": "MIT"},
			"test commit": {"demo": "this is test file"},
		})

		testutils.MakeFiles(t, fs, map[string]string{
			"untracked": "untracked content",
		})

		files, err := repo.GetFilesInIndex()

		require.NoError(t, err)
		assert.Empty(t, files)
	})

	t.Run("added files successfully", func(t *testing.T) {
		repo, _, fs, w := testutils.CreateRepo(t)

		testutils.MakeCommits(t, w, fs, map[string]map[string]string{
			"init commit": {"LICENSE": "MIT"},
			"test commit": {"demo": "this is test file"},
		})

		testutils.MakeFiles(t, fs, map[string]string{
			"tracked": "tracked content",
		})

		err := w.AddGlob(".")
		guards.NoError(err)

		files, err := repo.GetFilesInIndex()

		require.NoError(t, err)
		assert.Equal(t, []string{"tracked"}, files)
	})

	t.Run("worktree error", func(t *testing.T) {
		expectedErr := errors.New("worktree error")
		gitMock := mocks.NewGoGitRepositoryMock(t).WorktreeMock.Return(nil, expectedErr)

		repo := NewRepository(WithFactoryMethod(func() (GoGitRepository, storage.Storer, error) {
			return gitMock, nil, nil
		}))

		_, err := repo.GetFilesInIndex()

		require.EqualError(t, err, expectedErr.Error())
	})
}
