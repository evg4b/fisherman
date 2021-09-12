package vcs_test

import (
	"fisherman/pkg/guards"
	"fisherman/pkg/vcs"
	"fisherman/testing/testutils"
	"fmt"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
)

func TestGitRepository_GetCurrentBranch(t *testing.T) {
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

	assert.NoError(t, err)
	assert.Equal(t, expectedBranchName, branch)
}

func TestGitRepository_GetCurrentBranch_NoBranches(t *testing.T) {
	repo, _, fs, w := testutils.CreateRepo(t)

	testutils.MakeCommits(t, w, fs, map[string]map[string]string{
		"test commit": {"demo": "this is test file"},
	})

	branch, err := repo.GetCurrentBranch()

	assert.NoError(t, err)
	assert.Equal(t, "refs/heads/master", branch)
}

func TestGitRepository_GetCurrentBranch_NoCommits(t *testing.T) {
	repo, _, _, _ := testutils.CreateRepo(t)

	branch, err := repo.GetCurrentBranch()

	assert.NoError(t, err)
	assert.Equal(t, "", branch)
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

	assert.NoError(t, err)
	assert.Equal(t, vcs.User{
		UserName: expectedUserName,
		Email:    expectedEmail,
	}, user)
}
