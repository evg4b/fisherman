package vcs_test

import (
	"fisherman/internal/utils"
	"fisherman/testing/testutils"
	"fmt"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
)

func TestGitRepository_GetCurrentBranch(t *testing.T) {
	branchName := "test-branch"
	expectedBranchName := fmt.Sprintf("refs/heads/%s", branchName)

	repo, _, fs, w := testutils.CreateRepo()

	testutils.MakeCommits(w, fs, map[string]map[string]string{
		"test commit": {"demo": "this is test file"},
	})

	err := w.Checkout(&git.CheckoutOptions{
		Create: true,
		Branch: plumbing.NewBranchReferenceName(branchName),
	})
	utils.HandleCriticalError(err)

	branch, err := repo.GetCurrentBranch()

	assert.NoError(t, err)
	assert.Equal(t, expectedBranchName, branch)
}

func TestGitRepository_GetCurrentBranch_NoBranches(t *testing.T) {
	repo, _, fs, w := testutils.CreateRepo()

	testutils.MakeCommits(w, fs, map[string]map[string]string{
		"test commit": {"demo": "this is test file"},
	})

	branch, err := repo.GetCurrentBranch()

	assert.NoError(t, err)
	assert.Equal(t, "refs/heads/master", branch)
}
