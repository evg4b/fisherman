package testutils

import (
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func MakeCommits(wt *git.Worktree, fs billy.Basic, data map[string](map[string]string)) {
	for commitMessage, files := range data {
		MakeFiles(fs, files)
		err := wt.AddGlob(".")
		if err != nil {
			panic(err)
		}
		_, err = wt.Commit(commitMessage, &git.CommitOptions{
			Author: &object.Signature{
				Name:  "Test name",
				Email: "test@email.com",
			},
		})
		if err != nil {
			panic(err)
		}
	}
}

func MakeFiles(fs billy.Basic, files map[string]string) {
	for filemane, content := range files {
		err := util.WriteFile(fs, filemane, []byte(content), 0644)
		if err != nil {
			panic(err)
		}
	}
}
