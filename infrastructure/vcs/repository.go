package vcs

import (
	"fisherman/infrastructure"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

type GitRepository struct {
	repo *git.Repository
}

func NewGitRepository(path string) (*GitRepository, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	return &GitRepository{
		repo: r,
	}, nil
}

func (r *GitRepository) GetCurrentBranch() (string, error) {
	branchRefs, err := r.repo.Branches()
	if err != nil {
		return "", err
	}

	defer branchRefs.Close()

	headRef, err := r.repo.Head()
	if err != nil {
		return "", err
	}

	for branchRef, err := branchRefs.Next(); err == nil; {
		if branchRef.Hash() == headRef.Hash() {
			return branchRef.Name().String(), nil
		}
	}

	panic("Current branch not fount")
}

func (r *GitRepository) GetUser() (infrastructure.User, error) {
	config, err := r.repo.ConfigScoped(config.SystemScope)
	if err != nil {
		return infrastructure.User{}, err
	}

	return infrastructure.User{
		UserName: config.User.Name,
		Email:    config.User.Name,
	}, err
}
