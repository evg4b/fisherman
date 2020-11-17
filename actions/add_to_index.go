package actions

import (
	"errors"
	"fisherman/internal"

	"github.com/go-git/go-git/v5"
)

func AddToIndex(ctx internal.SyncContext, globs []string, optional bool) (bool, error) {
	if len(globs) > 0 {
		repo := ctx.Repository()
		for _, glob := range globs {
			err := repo.AddGlob(glob)
			if err != nil {
				if optional && errors.Is(err, git.ErrGlobNoMatches) {
					return true, nil
				}

				return false, err
			}
		}
	}

	return true, nil
}
