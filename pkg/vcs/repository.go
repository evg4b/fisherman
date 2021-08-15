package vcs

import (
	"errors"
	"fisherman/internal"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

type GitRepository struct {
	repo func() (*git.Repository, error)
}

func OpenGitRepository(path string) *GitRepository {
	var repoOnce sync.Once
	var repo *git.Repository

	return &GitRepository{
		repo: func() (*git.Repository, error) {
			var err error

			repoOnce.Do(func() {
				repo, err = git.PlainOpen(path)
			})

			return repo, err
		},
	}
}

func CreateGitRepository(repo *git.Repository) *GitRepository {
	return &GitRepository{
		repo: func() (*git.Repository, error) {
			return repo, nil
		},
	}
}

func (r *GitRepository) GetCurrentBranch() (string, error) {
	repo, err := r.repo()
	if err != nil {
		return "", err
	}

	headRef, err := repo.Head()
	if err != nil {
		if errors.Is(err, plumbing.ErrReferenceNotFound) {
			return "", nil
		}

		return "", err
	}

	return headRef.Name().String(), nil
}

func (r *GitRepository) GetUser() (internal.User, error) {
	repo, err := r.repo()
	if err != nil {
		return internal.User{}, err
	}

	gitConfig, err := repo.ConfigScoped(config.SystemScope)
	if err != nil {
		return internal.User{}, err
	}

	return internal.User{
		UserName: gitConfig.User.Name,
		Email:    gitConfig.User.Name,
	}, err
}

func (r *GitRepository) AddGlob(glob string) error {
	repo, err := r.repo()
	if err != nil {
		return err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	return wt.AddGlob(glob)
}

func (r *GitRepository) RemoveGlob(glob string) error {
	repo, err := r.repo()
	if err != nil {
		return err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	return wt.RemoveGlob(glob)
}
