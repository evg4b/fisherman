package vcs

import (
	"errors"
	"fisherman/infrastructure"
	"fisherman/utils"
	"io"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type GitRepository struct {
	path         string
	internalRepo *git.Repository
}

func NewGitRepository(path string) (*GitRepository, error) {
	return &GitRepository{path: path, internalRepo: nil}, nil
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

func (r *GitRepository) GetLastTag() (string, error) {
	tagRef, err := r.repo().Tags()
	if err != nil {
		if errors.Is(err, plumbing.ErrReferenceNotFound) {
			return "", nil
		}

		return "", err
	}

	var latestTagCommit *object.Commit
	var latestTagName string

	defer tagRef.Close()
	for {
		tagRef, err := tagRef.Next()
		if err != nil {
			if err == io.EOF {
				return latestTagName, nil
			}

			return "", err
		}

		revision := plumbing.Revision(tagRef.Name().String())
		tagCommitHash, err := r.repo().ResolveRevision(revision)
		if err != nil {
			return "", err
		}

		commit, err := r.repo().CommitObject(*tagCommitHash)
		if err != nil {
			return "", err
		}

		if latestTagCommit == nil {
			latestTagCommit = commit
			latestTagName = tagRef.Name().String()
		}

		if commit.Committer.When.After(latestTagCommit.Committer.When) {
			latestTagCommit = commit
			latestTagName = tagRef.Name().String()
		}
	}
}

func (r *GitRepository) repo() *git.Repository {
	if r.internalRepo == nil {
		repo, err := git.PlainOpen(r.path)
		utils.HandleCriticalError(err)
		r.internalRepo = repo
	}

	return r.internalRepo
}
