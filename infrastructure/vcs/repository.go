package vcs

import (
	"errors"
	"fisherman/infrastructure"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
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
	headRef, err := r.repo.Head()
	if err != nil {
		if errors.Is(err, plumbing.ErrReferenceNotFound) {
			return "", nil
		}

		return "", err
	}

	return headRef.Name().String(), nil
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
