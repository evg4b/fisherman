package vcs

import (
	"github.com/go-git/go-git/v5"
)

func (r *GitRepository) GetFilesInIndex() ([]string, error) {
	worktree, err := r.repo().Worktree()
	if err != nil {
		return nil, err
	}

	statusIndex, err := worktree.Status()
	if err != nil {
		return nil, err
	}

	files := []string{}
	for key, status := range statusIndex {
		if isModified(status.Staging) {
			files = append(files, key)
		}
	}

	return files, nil
}

func isModified(status git.StatusCode) bool {
	return status != git.Unmodified && status != git.Untracked
}
