package vcs

import (
	"errors"
	"io"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// nolint: cyclop
func (r *GitRepository) GetLastTag() (string, error) {
	repo, err := r.repo()
	if err != nil {
		return "", err
	}

	tagRef, err := repo.Tags()
	if err != nil {
		if errors.Is(err, plumbing.ErrReferenceNotFound) {
			return "", nil
		}

		return "", err
	}

	head, err := repo.Head()
	if err != nil {
		if errors.Is(err, plumbing.ErrReferenceNotFound) {
			return "", nil
		}

		return "", err
	}

	headCommit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return "", err
	}

	var latestTagCommit *object.Commit
	var latestTagName string

	defer tagRef.Close()
	for {
		tagRef, err := tagRef.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return latestTagName, nil
			}

			return "", err
		}

		revision := plumbing.Revision(tagRef.Name().String())
		tagCommitHash, err := repo.ResolveRevision(revision)
		if err != nil {
			return "", err
		}

		commit, err := repo.CommitObject(*tagCommitHash)
		if err != nil {
			return "", err
		}

		if commit.Committer.When.After(headCommit.Committer.When) {
			continue
		}

		if latestTagCommit == nil {
			latestTagCommit = commit
			latestTagName = tagRef.Name().String()
		}

		if commit.Committer.When.After(latestTagCommit.Committer.When) || commit.Committer.When.Equal(latestTagCommit.Committer.When) {
			latestTagCommit = commit
			latestTagName = tagRef.Name().String()
		}
	}
}
