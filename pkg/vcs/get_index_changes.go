package vcs

import (
	"bytes"
	"io/ioutil"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/utils/diff"
	"github.com/go-git/go-git/v5/utils/merkletrie"
	"github.com/go-git/go-git/v5/utils/merkletrie/index"
	"github.com/go-git/go-git/v5/utils/merkletrie/noder"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func (r *GitRepository) GetIndexChanges() (map[string]Changes, error) {
	indexChanges := make(map[string]Changes)

	repo, err := r.repo()
	if err != nil {
		return nil, err
	}

	head, err := repo.Head()
	if err != nil {
		return nil, err
	}

	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return nil, err
	}

	commitTree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	from := object.NewTreeRootNode(commitTree)

	idx, err := repo.Storer.Index()
	if err != nil {
		return nil, err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	fS := wt.Filesystem

	to := index.NewRootNode(idx)

	dd, err := merkletrie.DiffTree(from, to, diffTreeIsEquals)
	if err != nil {
		return nil, err
	}

	for _, v := range dd {
		ff, _ := fS.Open(v.To.Name())
		content, err := ioutil.ReadAll(ff)
		if err != nil {
			return nil, err
		}

		ff.Close()

		if v.From == nil {

			indexChanges[v.To.Name()] = Changes{
				Change{
					Status: Added,
					Change: string(content),
				},
			}
		} else {

			ff, err := commitTree.File("tracked")
			if err != nil {
				return nil, err
			}

			cc, _ := ff.Contents()

			difference := diff.Do(cc, string(content))

			file := Changes{}
			for _, ddd := range difference {
				if ddd.Type != diffmatchpatch.DiffEqual {
					file = append(file, Change{
						Status: convertStatis(ddd.Type),
						Change: ddd.Text,
					})
				}
			}

			indexChanges[v.To.Name()] = file
		}
	}

	return indexChanges, nil
}

func convertStatis(ty diffmatchpatch.Operation) ChangeCode {
	switch ty {
	case diffmatchpatch.DiffEqual:
		return Unmodified
	case diffmatchpatch.DiffInsert:
		return Added
	case diffmatchpatch.DiffDelete:
		return Deleted
	}

	panic("sadsad")
}

func isModified(status git.StatusCode) bool {
	return status != git.Unmodified && status != git.Untracked
}

var emptyNoderHash = make([]byte, 24)

func diffTreeIsEquals(a, b noder.Hasher) bool {
	hashA := a.Hash()
	hashB := b.Hash()

	if bytes.Equal(hashA, emptyNoderHash) || bytes.Equal(hashB, emptyNoderHash) {
		return false
	}

	return bytes.Equal(hashA, hashB)
}
