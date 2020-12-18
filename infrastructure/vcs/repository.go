package vcs

import (
	"errors"
	"fisherman/infrastructure"
	"fisherman/utils"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

type GitRepository struct {
	path         string
	internalRepo *git.Repository
}

func NewGitRepository(path string) *GitRepository {
	return &GitRepository{path: path, internalRepo: nil}
}

func (r *GitRepository) GetCurrentBranch() (string, error) {
	headRef, err := r.repo().Head()
	if err != nil {
		if errors.Is(err, plumbing.ErrReferenceNotFound) {
			return "", nil
		}

		return "", err
	}

	return headRef.Name().String(), nil
}

func (r *GitRepository) GetUser() (infrastructure.User, error) {
	gitConfig, err := r.repo().ConfigScoped(config.SystemScope)
	if err != nil {
		return infrastructure.User{}, err
	}

	return infrastructure.User{
		UserName: gitConfig.User.Name,
		Email:    gitConfig.User.Name,
	}, err
}

func (r *GitRepository) repo() *git.Repository {
	if r.internalRepo == nil {
		repo, err := git.PlainOpen(r.path)
		utils.HandleCriticalError(err)
		r.internalRepo = repo
	}

	return r.internalRepo
}

func (r *GitRepository) AddGlob(glob string) error {
	wt, err := r.repo().Worktree()
	if err != nil {
		return err
	}

	return wt.AddGlob(glob)
}

func (r *GitRepository) RemoveGlob(glob string) error {
	wt, err := r.repo().Worktree()
	if err != nil {
		return err
	}

	return wt.RemoveGlob(glob)
}
