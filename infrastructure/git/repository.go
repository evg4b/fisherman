package git

import (
	"fisherman/utils"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type GitRepository struct {
	repo *git.Repository
}

func NewRepository(path string) *GitRepository {
	r, err := git.PlainOpen(path)
	utils.HandleCriticalError(err)

	return &GitRepository{
		repo: r,
	}
}

func (r *GitRepository) GetCurrentBranch() (string, error) {
	branchRefs, err := r.repo.Branches()
	if err != nil {
		return "", err
	}

	headRef, err := r.repo.Head()
	if err != nil {
		return "", err
	}

	var currentBranchName string
	err = branchRefs.ForEach(func(branchRef *plumbing.Reference) error {
		if branchRef.Hash() == headRef.Hash() {
			currentBranchName = branchRef.Name().String()

			return nil
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return currentBranchName, nil
}
