package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"path/filepath"
)

type RepositoryInfo struct {
	Path          string
	CurrentBranch string
}

func GetRepositoryInfo(path string) (*RepositoryInfo, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	r, err := git.PlainOpen(absPath)
	if err != nil {
		return nil, err
	}

	branch, err := getCurrentBranchFromRepository(r)
	if err != nil {
		return nil, err
	}

	return &RepositoryInfo{
		Path:          absPath,
		CurrentBranch: branch,
	}, nil
}

func getCurrentBranchFromRepository(repository *git.Repository) (string, error) {
	branchRefs, err := repository.Branches()
	if err != nil {
		return "", err
	}

	headRef, err := repository.Head()
	if err != nil {
		return "", err
	}

	var currentBranchName string
	err = branchRefs.ForEach(func(branchRef *plumbing.Reference) error {
		if branchRef.Hash() == headRef.Hash() {
			currentBranchName = branchRef.Name().String()
			return nil
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return currentBranchName, nil
}
