package vcs_test

import (
	"errors"
	"fmt"
	"github.com/evg4b/fisherman/pkg/guards"
	"github.com/evg4b/fisherman/testing/testutils"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/evg4b/fisherman/pkg/vcs"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage"
	"github.com/stretchr/testify/assert"
)

func TestGitRepository_GetCurrentBranch(t *testing.T) {
	t.Run("returns currect branch name", func(t *testing.T) {
		branchName := "test-branch"
		expectedBranchName := fmt.Sprintf("refs/heads/%s", branchName)

		repo, _, fs, w := testutils.CreateRepo(t)

		testutils.MakeCommits(t, w, fs, map[string]map[string]string{
			"test commit": {"demo": "this is test file"},
		})

		err := w.Checkout(&git.CheckoutOptions{
			Create: true,
			Branch: plumbing.NewBranchReferenceName(branchName),
		})
		guards.NoError(err)

		branch, err := repo.GetCurrentBranch()

		require.NoError(t, err)
		assert.Equal(t, expectedBranchName, branch)
	})

	t.Run("no branches", func(t *testing.T) {
		repo, _, fs, w := testutils.CreateRepo(t)

		testutils.MakeCommits(t, w, fs, map[string]map[string]string{
			"test commit": {"demo": "this is test file"},
		})

		branch, err := repo.GetCurrentBranch()

		require.NoError(t, err)
		assert.Equal(t, "refs/heads/master", branch)
	})

	t.Run("no commits in branch", func(t *testing.T) {
		repo, _, _, _ := testutils.CreateRepo(t)

		branch, err := repo.GetCurrentBranch()

		require.NoError(t, err)
		assert.Equal(t, "", branch)
	})
}

func TestGitRepository_GetUser(t *testing.T) {
	expectedUserName := "TestUser"
	expectedEmail := "TestUser@mail.com"

	repo, r, _, _ := testutils.CreateRepo(t)

	err := r.SetConfig(&config.Config{
		User: struct {
			Name  string
			Email string
		}{
			Name:  expectedUserName,
			Email: expectedEmail,
		},
	})
	guards.NoError(err)

	user, err := repo.GetUser()

	require.NoError(t, err)
	assert.Equal(t, User{
		UserName: expectedUserName,
		Email:    expectedEmail,
	}, user)
}

func TestLazyInitialization(t *testing.T) {
	test := []struct {
		name   string
		action func(repo *GitRepository) error
	}{
		{
			name: "AddGlob initialize git repository",
			action: func(repo *GitRepository) error {
				return repo.AddGlob("demo.ts")
			},
		},
		{
			name: "AddRemoveGlobGlob initialize git repository",
			action: func(repo *GitRepository) error {
				return repo.RemoveGlob("demo.ts")
			},
		},
		{
			name: "GetUser initialize git repository",
			action: func(repo *GitRepository) error {
				_, err := repo.GetUser()

				return err
			},
		},
		{
			name: "GetCurrentBranch initialize git repository",
			action: func(repo *GitRepository) error {
				_, err := repo.GetCurrentBranch()

				return err
			},
		},
		{
			name: "GetLastTag initialize git repository",
			action: func(repo *GitRepository) error {
				_, err := repo.GetLastTag()

				return err
			},
		},
		{
			name: "GetIndexChanges initialize git repository",
			action: func(repo *GitRepository) error {
				_, err := repo.GetIndexChanges()

				return err
			},
		},
		{
			name: "GetFilesInIndex initialize git repository",
			action: func(repo *GitRepository) error {
				_, err := repo.GetFilesInIndex()

				return err
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewRepository(WithFactoryMethod(func() (GoGitRepository, storage.Storer, error) {
				return nil, nil, errors.New("Initialization error")
			}))

			err := tt.action(repo)

			require.Error(t, err, "Initialization error")
		})
	}
}
