package testutils

import (
	"fisherman/pkg/guards"
	"testing"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var index int = 0

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

func MakeFiles(t *testing.T, fs billy.Basic, files map[string]string) {
	t.Helper()

	for filemane, content := range files {
		err := util.WriteFile(fs, filemane, []byte(content), 0644)
		guards.NoError(err)
	}
}
