package git

import (
	"fisherman/utils"

	"github.com/go-git/go-git/v5"
)

type FGitRepository struct {
	repo *git.Repository
}

func NewRepository(path string) *FGitRepository {
	r, err := git.PlainOpen(path)
	utils.HandleCriticalError(err)

	return &FGitRepository{
		repo: r,
	}
}

func (r *FGitRepository) GetCurrentBranch() (string, error) {
	branchRefs, err := r.repo.Branches()
	if err != nil {
		return "", err
	}

	defer branchRefs.Close()

	headRef, err := r.repo.Head()
	if err != nil {
		return "", err
	}

	for branchRef, err := branchRefs.Next(); err != nil; {
		if branchRef.Hash() == headRef.Hash() {
			return branchRef.Name().String(), nil
		}
	}

	panic("Current branch not fount")
}
