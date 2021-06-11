package rules

import (
	"fisherman/internal"
	"fisherman/internal/utils"
	"io"
	"path/filepath"

	"github.com/go-errors/errors"

	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/go-multierror"
)

const SuppressCommitFilesType = "suppress-commit-files"

type SuppressCommitFiles struct {
	BaseRule        `yaml:",inline"`
	Globs           []string `yaml:"globs"`
	RemoveFromIndex bool     `yaml:"remove-from-index"`
}

func (rule *SuppressCommitFiles) Check(ctx internal.ExecutionContext, _ io.Writer) error {
	if len(rule.Globs) == 0 {
		return nil
	}

	repo := ctx.Repository()
	files, err := repo.GetFilesInIndex()
	if err != nil {
		return err
	}

	matchedFiles := []string{}
	for _, glob := range rule.Globs {
		for _, file := range files {
			matched, err := filepath.Match(glob, file)
			if err != nil {
				return err
			}
			if matched {
				matchedFiles = append(matchedFiles, file)
			}
		}
	}

	var multiError *multierror.Error
	for _, file := range matchedFiles {
		if rule.RemoveFromIndex {
			err := repo.RemoveGlob(file)
			if err != nil && !errors.Is(err, git.ErrGlobNoMatches) {
				return err
			}
		} else {
			multiError = multierror.Append(multiError, errors.Errorf("file %s can not be committed", file))
		}
	}

	return multiError.ErrorOrNil()
}

func (rule *SuppressCommitFiles) GetPosition() byte {
	return PostScripts
}

func (rule *SuppressCommitFiles) Compile(variables map[string]interface{}) {
	rule.BaseRule.Compile(variables)
	utils.FillTemplatesArray(rule.Globs, variables)
}
