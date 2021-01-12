package rules

import (
	"errors"
	"fisherman/internal"
	"fmt"
	"io"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/go-multierror"
)

const SuppressCommitType = "suppress-commit-files"

type SuppressCommitFiles struct {
	BaseRule        `mapstructure:",squash"`
	Globs           []string `mapstructure:"globs"`
	RemoveFromIndex bool     `mapstructure:"remove-from-index"`
}

func (rule SuppressCommitFiles) Check(_ io.Writer, ctx internal.AsyncContext) error {
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
			multiError = multierror.Append(multiError, fmt.Errorf("file %s can not be committed", file))
		}
	}

	return multiError.ErrorOrNil()
}

func (rule SuppressCommitFiles) GetPosition() byte {
	return AfterScripts
}
