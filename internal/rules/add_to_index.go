package rules

import (
	"errors"
	"fisherman/internal"
	"io"

	"github.com/go-git/go-git/v5"
)

const AddToIndexType = "add-to-index"

type Glob struct {
	Glob       string `mapstructure:"glob"`
	IsRequired bool   `mapstructure:"required"`
}

type AddToIndex struct {
	BaseRule `mapstructure:",squash"`
	Globs    []Glob `mapstructure:"globs"`
}

func (rule *AddToIndex) GetPosition() byte {
	return PostScripts
}

func (rule AddToIndex) Check(ctx internal.ExecutionContext, _ io.Writer) error {
	if len(rule.Globs) > 0 {
		repo := ctx.Repository()
		for _, glob := range rule.Globs {
			err := repo.AddGlob(glob.Glob)
			if err != nil {
				if errors.Is(err, git.ErrGlobNoMatches) && !glob.IsRequired {
					continue
				}

				return err
			}
		}
	}

	return nil
}
