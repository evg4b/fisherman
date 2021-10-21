package vcs

import (
	"errors"

	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage"
)

type GitRepository struct {
	init   func() error
	repo   GoGitRepository
	storer storage.Storer
}

// NewRepository returns not initialized git repo with passed options.
//
// Repo has lazy initialization. Error in case of error opening the repository
// will be returned only on the first access to any method.
func NewRepository(options ...repositoryOption) *GitRepository {
	repo := &GitRepository{}

	for _, option := range options {
		option(repo)
	}

	return repo
}

// GetCurrentBranch return current branch name.
func (r *GitRepository) GetCurrentBranch() (string, error) {
	if err := r.init(); err != nil {
		return "", err
	}

	headRef, err := r.repo.Head()
	if err != nil {
		if errors.Is(err, plumbing.ErrReferenceNotFound) {
			return "", nil
		}

		return "", err
	}

	return headRef.Name().String(), nil
}

// GetUser return information ablout configured git user.
func (r *GitRepository) GetUser() (User, error) {
	if err := r.init(); err != nil {
		return User{}, err
	}

	gitConfig, err := r.repo.ConfigScoped(config.SystemScope)
	if err != nil {
		return User{}, err
	}

	return User{
		UserName: gitConfig.User.Name,
		Email:    gitConfig.User.Email,
	}, nil
}

// AddGlob adds files in index by glob expresion.
func (r *GitRepository) AddGlob(glob string) error {
	if err := r.init(); err != nil {
		return err
	}

	wt, err := r.repo.Worktree()
	if err != nil {
		return err
	}

	return wt.AddGlob(glob)
}

// AddGlob removes files from index by glob expresion.
func (r *GitRepository) RemoveGlob(glob string) error {
	if err := r.init(); err != nil {
		return err
	}

	wt, err := r.repo.Worktree()
	if err != nil {
		return err
	}

	return wt.RemoveGlob(glob)
}
