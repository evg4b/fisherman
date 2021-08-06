package testutils

import (
	"fisherman/internal/utils"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var index int = 0

func MakeCommits(wt *git.Worktree, fs billy.Basic, data map[string](map[string]string)) {
	for commitMessage, files := range data {
		MakeFiles(fs, files)
		err := wt.AddGlob(".")
		utils.HandleCriticalError(err)

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
		utils.HandleCriticalError(err)
	}
}

func MakeFiles(fs billy.Basic, files map[string]string) {
	for filemane, content := range files {
		err := util.WriteFile(fs, filemane, []byte(content), 0644)
		utils.HandleCriticalError(err)
	}
}