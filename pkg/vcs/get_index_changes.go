package vcs

import (
	"bytes"
	"path"

	"github.com/go-git/go-billy/v5/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/utils/diff"
	"github.com/go-git/go-git/v5/utils/merkletrie"
	"github.com/go-git/go-git/v5/utils/merkletrie/index"
	"github.com/go-git/go-git/v5/utils/merkletrie/noder"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// nolint: cyclop
func (r *GitRepository) GetIndexChanges() (map[string]Changes, error) {
	indexChanges := make(map[string]Changes)

	if err := r.init(); err != nil {
		return nil, err
	}

	head, err := r.repo.Head()
	if err != nil {
		return nil, err
	}

	commit, err := r.repo.CommitObject(head.Hash())
	if err != nil {
		return nil, err
	}

	commitTree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	idx, err := r.storer.Index()
	if err != nil {
		return nil, err
	}

	wt, err := r.repo.Worktree()
	if err != nil {
		return nil, err
	}

	fs := wt.Filesystem

	diffTree, err := merkletrie.DiffTree(object.NewTreeRootNode(commitTree), index.NewRootNode(idx), diffTreeIsEquals)
	if err != nil {
		return nil, err
	}

	for _, diffTreeItem := range diffTree {
		toPath := convertToPath(diffTreeItem.To)
		toContentBytes, err := util.ReadFile(fs, toPath)
		if err != nil {
			return nil, err
		}

		toContent := string(toContentBytes)

		if diffTreeItem.From == nil {
			indexChanges[toPath] = Changes{
				Change{Status: Added, Change: toContent},
			}

			continue
		}

		fromPath := convertToPath(diffTreeItem.From)
		fromFile, err := commitTree.File(fromPath)
		if err != nil {
			return nil, err
		}

		fromContent, err := fromFile.Contents()
		if err != nil {
			return nil, err
		}

		fileChanges := Changes{}
		for _, diffItem := range diff.Do(fromContent, toContent) {
			if diffItem.Type != diffmatchpatch.DiffEqual {
				fileChanges = append(fileChanges, Change{
					Status: convertStatus(diffItem.Type),
					Change: diffItem.Text,
				})
			}
		}

		indexChanges[toPath] = fileChanges
	}

	return indexChanges, nil
}

func convertToPath(node noder.Path) string {
	pathValue := make([]string, 0, len(node))

	for _, nodeSection := range node {
		pathValue = append(pathValue, nodeSection.Name())
	}

	return path.Join(pathValue...)
}

func convertStatus(ty diffmatchpatch.Operation) ChangeCode {
	switch ty {
	case diffmatchpatch.DiffEqual:
		return Unmodified
	case diffmatchpatch.DiffInsert:
		return Added
	case diffmatchpatch.DiffDelete:
		return Deleted
	}

	panic("incorrect diffmatchpatch.Operation")
}

func isModified(status git.StatusCode) bool {
	return status != git.Unmodified && status != git.Untracked
}

const hashSize = 24

var emptyHash = make([]byte, hashSize)

func diffTreeIsEquals(a, b noder.Hasher) bool {
	hashA := a.Hash()
	hashB := b.Hash()

	if bytes.Equal(hashA, hashB) && !bytes.Equal(hashA, emptyHash) {
		return true
	}

	return false
}
