package actions

import (
	"errors"
	"fisherman/internal"

	"github.com/go-git/go-git/v5"
)

type Glob struct {
	Glob       string `yaml:"glob"`
	IsRequired bool   `yaml:"required"`
}

func (glob *Glob) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var stringGlob string
	err := unmarshal(&stringGlob)
	if err == nil {
		glob.Glob = stringGlob
		glob.IsRequired = true

		return nil
	}

	type plain Glob

	return unmarshal((*plain)(glob))
}

func AddToIndex(ctx internal.AsyncContext, globs []Glob) (bool, error) {
	if len(globs) > 0 {
		repo := ctx.Repository()
		for _, glob := range globs {
			err := repo.AddGlob(glob.Glob)
			if err != nil {
				if errors.Is(err, git.ErrGlobNoMatches) && !glob.IsRequired {
					continue
				}

				return false, err
			}
		}
	}

	return true, nil
}
