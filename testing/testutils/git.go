package testutils

import (
	"testing"
	"time"

	"fisherman/pkg/guards"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var index int

// MakeCommits creates commits history from two-dimensional map.
// First level key is a commit message, second key is a file name and value is a file content.
func MakeCommits(t *testing.T, wt *git.Worktree, fs billy.Basic, data map[string](map[string]string)) {
	t.Helper()

	for commitMessage, files := range data {
		MakeFiles(t, fs, files)
		err := wt.AddGlob(".")
		guards.NoError(err)

		signature := object.Signature{
			Name:  "Test name",
			Email: "test@email.com",
			When:  time.Now().Add(time.Minute * time.Duration(index)),
		}
		_, err = wt.Commit(commitMessage, &git.CommitOptions{
			Author:    &signature,
			Committer: &signature,
		})
		index++
		guards.NoError(err)
	}
}
