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
	path         string
	internalRepo *git.Repository
	repoOnce     sync.Once
}

func NewGitRepository(path string) *GitRepository {
	return &GitRepository{path: path, internalRepo: nil}
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

func (r *GitRepository) repo() (*git.Repository, error) {
	var err error

	r.repoOnce.Do(func() {
		r.internalRepo, err = git.PlainOpen(r.path)
	})

	return r.internalRepo, err
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
