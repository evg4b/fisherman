package rules

import (
	"errors"
	"fisherman/internal"
	"fisherman/internal/utils"
	"io"

	"github.com/go-git/go-git/v5"
)

const AddToIndexType = "add-to-index"

type Glob struct {
	Glob       string `yaml:"glob"`
	IsRequired bool   `yaml:"required"`
}

type AddToIndex struct {
	BaseRule `yaml:",inline"`
	Globs    []Glob `yaml:"globs"`
}

func (rule *AddToIndex) GetPosition() byte {
	return PostScripts
}

func (rule *AddToIndex) Check(ctx internal.ExecutionContext, _ io.Writer) error {
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

func (rule *AddToIndex) Compile(variables map[string]interface{}) {
	rule.BaseRule.Compile(variables)
	for index := range rule.Globs {
		utils.FillTemplate(&rule.Globs[index].Glob, variables)
	}
}
