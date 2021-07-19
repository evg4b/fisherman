package vcs_test

import (
	"fisherman/internal/utils"
	"fisherman/testing/testutils"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
)

func TestGitRepository_GetLastTag(t *testing.T) {
	repo, r, fs, w := createRepo()

	testutils.MakeCommits(w, fs, map[string]map[string]string{
		"init commit": {"LICENSE": "MIT"},
	})

	head, err := r.Head()
	utils.HandleCriticalError(err)

	_, err = r.CreateTag("tag1", head.Hash(), &git.CreateTagOptions{
		Message: "test tag 1",
		Tagger: &object.Signature{
			Name:  "Test name",
			Email: "test@email.com",
		},
	})
	utils.HandleCriticalError(err)

	testutils.MakeCommits(w, fs, map[string]map[string]string{
		"test commit": {"demo": "this is test file"},
	})

	head, err = r.Head()
	utils.HandleCriticalError(err)

	_, err = r.CreateTag("tag2", head.Hash(), &git.CreateTagOptions{
		Message: "test tag 2",
		Tagger: &object.Signature{
			Name:  "Test name",
			Email: "test@email.com",
		},
	})
	utils.HandleCriticalError(err)

	tag, err := repo.GetLastTag()

	assert.NoError(t, err)
	assert.Equal(t, "refs/tags/tag2", tag)
}

func TestGitRepository_GetLastTag_NotLastHead(t *testing.T) {
	repo, r, fs, w := createRepo()

	testutils.MakeCommits(w, fs, map[string]map[string]string{
		"init commit": {"LICENSE": "MIT"},
	})

	expectedCommitRef, err := r.Head()
	utils.HandleCriticalError(err)

	_, err = r.CreateTag("tag1", expectedCommitRef.Hash(), &git.CreateTagOptions{
		Message: "test tag 1",
		Tagger: &object.Signature{
			Name:  "Test name",
			Email: "test@email.com",
		},
	})
	utils.HandleCriticalError(err)

	testutils.MakeCommits(w, fs, map[string]map[string]string{
		"test commit": {"demo": "this is test file"},
	})

	head, err := r.Head()
	utils.HandleCriticalError(err)

	_, err = r.CreateTag("tag2", head.Hash(), &git.CreateTagOptions{
		Message: "test tag 2",
		Tagger: &object.Signature{
			Name:  "Test name",
			Email: "test@email.com",
		},
	})
	utils.HandleCriticalError(err)

	err = w.Checkout(&git.CheckoutOptions{
		Hash: expectedCommitRef.Hash(),
	})
	utils.HandleCriticalError(err)

	tag, err := repo.GetLastTag()

	assert.NoError(t, err)
	assert.Equal(t, "refs/tags/tag1", tag)
}

func TestGitRepository_GetLastTag_NoTags(t *testing.T) {
	repo, _, fs, w := createRepo()

	testutils.MakeCommits(w, fs, map[string]map[string]string{
		"init commit": {"LICENSE": "MIT"},
		"test commit": {"demo": "this is test file"},
	})

	tag, err := repo.GetLastTag()

	assert.NoError(t, err)
	assert.Empty(t, tag)
}

func TestGitRepository_GetLastTag_EmptyRepo(t *testing.T) {
	repo, _, _, _ := createRepo()

	tag, err := repo.GetLastTag()

	assert.NoError(t, err)
	assert.Empty(t, tag)
}
