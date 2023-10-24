package vcs_test

import (
	"github.com/evg4b/fisherman/pkg/guards"
	"github.com/evg4b/fisherman/testing/testutils"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
)

func TestGitRepository_GetLastTag(t *testing.T) {
	t.Run("head is last commit", func(t *testing.T) {
		repo, r, fs, w := testutils.CreateRepo(t)

		testutils.MakeCommits(t, w, fs, map[string]map[string]string{
			"init commit": {"LICENSE": "MIT"},
		})

		head, err := r.Head()
		guards.NoError(err)

		_, err = r.CreateTag("tag1", head.Hash(), &git.CreateTagOptions{
			Message: "test tag 1",
			Tagger: &object.Signature{
				Name:  "Test name",
				Email: "test@email.com",
			},
		})
		guards.NoError(err)

		testutils.MakeCommits(t, w, fs, map[string]map[string]string{
			"test commit": {"demo": "this is test file"},
		})

		head, err = r.Head()
		guards.NoError(err)

		_, err = r.CreateTag("tag2", head.Hash(), &git.CreateTagOptions{
			Message: "test tag 2",
			Tagger: &object.Signature{
				Name:  "Test name",
				Email: "test@email.com",
			},
		})
		guards.NoError(err)

		tag, err := repo.GetLastTag()

		require.NoError(t, err)
		assert.Equal(t, "refs/tags/tag2", tag)
	})

	t.Run("head is not last commint", func(t *testing.T) {
		repo, r, fs, w := testutils.CreateRepo(t)

		testutils.MakeCommits(t, w, fs, map[string]map[string]string{
			"init commit": {"LICENSE": "MIT"},
		})

		expectedCommitRef, err := r.Head()
		guards.NoError(err)

		_, err = r.CreateTag("tag1", expectedCommitRef.Hash(), &git.CreateTagOptions{
			Message: "test tag 1",
			Tagger: &object.Signature{
				Name:  "Test name",
				Email: "test@email.com",
			},
		})
		guards.NoError(err)

		testutils.MakeCommits(t, w, fs, map[string]map[string]string{
			"test commit": {"demo": "this is test file"},
		})

		head, err := r.Head()
		guards.NoError(err)

		_, err = r.CreateTag("tag2", head.Hash(), &git.CreateTagOptions{
			Message: "test tag 2",
			Tagger: &object.Signature{
				Name:  "Test name",
				Email: "test@email.com",
			},
		})
		guards.NoError(err)

		err = w.Checkout(&git.CheckoutOptions{
			Hash: expectedCommitRef.Hash(),
		})
		guards.NoError(err)

		tag, err := repo.GetLastTag()

		require.NoError(t, err)
		assert.Equal(t, "refs/tags/tag1", tag)
	})

	t.Run("no commits", func(t *testing.T) {
		repo, _, _, _ := testutils.CreateRepo(t)

		tag, err := repo.GetLastTag()

		require.NoError(t, err)
		assert.Equal(t, "", tag)
	})

	t.Run("no tags", func(t *testing.T) {
		repo, _, fs, w := testutils.CreateRepo(t)

		testutils.MakeCommits(t, w, fs, map[string]map[string]string{
			"init commit": {"LICENSE": "MIT"},
			"test commit": {"demo": "this is test file"},
		})

		tag, err := repo.GetLastTag()

		require.NoError(t, err)
		assert.Empty(t, tag)
	})

	t.Run("empty repository", func(t *testing.T) {
		repo, _, _, _ := testutils.CreateRepo(t)

		tag, err := repo.GetLastTag()

		require.NoError(t, err)
		assert.Empty(t, tag)
	})
}
