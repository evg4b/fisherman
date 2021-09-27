package vcs

import (
	"errors"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

type GitRepository struct {
	repo func() (*git.Repository, error)
}

// OpenGitRepository returns not initialized git repo with root in passed path.
//
// Repo has lazy initialization. Error in case of error opening the repository
// will be returned only on the first access to any method.
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

// CreateGitRepository returns not initialized git repo based on passed go-git repo
//
// Repo has lazy initialization. Error in case of error opening the repository
// will be returned only on the first access to any method.
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

func (r *GitRepository) GetUser() (User, error) {
	repo, err := r.repo()
	if err != nil {
		return User{}, err
	}

	gitConfig, err := repo.ConfigScoped(config.SystemScope)
	if err != nil {
		return User{}, err
	}

	return User{
		UserName: gitConfig.User.Name,
		Email:    gitConfig.User.Email,
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
