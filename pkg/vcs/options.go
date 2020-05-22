package vcs

import (
	"sync"

	"github.com/go-git/go-git/v5"
)

func WithFsRepo(path string) RepositoryOption {
	var repoOnce sync.Once
	var gitRepo *git.Repository

	return func(repo *GitRepository) {
		repo.init = func() error {
			var err error

			repoOnce.Do(func() {
				gitRepo, err = git.PlainOpen(path)
				repo.repo = gitRepo
				repo.storer = gitRepo.Storer
			})

			return err
		}
	}
}

func WithFactoryMethod(factory factoryMethod) RepositoryOption {
	var repoOnce sync.Once

	return func(repo *GitRepository) {
		repo.init = func() error {
			var err error

			repoOnce.Do(func() {
				repo.repo, repo.storer, err = factory()
			})

			return err
		}
	}
}

func WithGitRepository(gitRepo *git.Repository) RepositoryOption {
	return func(repo *GitRepository) {
		repo.init = func() error {
			return nil
		}
		repo.repo = gitRepo
		repo.storer = gitRepo.Storer
	}
}
