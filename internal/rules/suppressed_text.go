package rules

import (
	"context"
	"fmt"
	"github.com/evg4b/fisherman/internal/utils"
	"io"
	"strings"

	"github.com/hashicorp/go-multierror"
)

const SuppressedTextType = "suppressed-text"

type SuppressedText struct {
	BaseRule      `yaml:",inline"`
	Substrings    []string `yaml:"substrings"`
	ExcludedGlobs []string `yaml:"exclude"`
}

func (rule *SuppressedText) Check(_ context.Context, _ io.Writer) error {
	repo := rule.BaseRule.repo
	changes, err := repo.GetIndexChanges()
	if err != nil {
		return err
	}

	var validationErrors multierror.Error
	validationErrors.ErrorFormat = plainErrorFormatter

	for file, fileChanges := range changes {
		matched, err := utils.MatchToGlobs(rule.ExcludedGlobs, file)
		if err != nil {
			return err
		}

		if matched {
			continue
		}

		addedLines := fileChanges.Added()
		for _, suppressed := range rule.Substrings {
			for _, lineChange := range addedLines {
				if strings.Contains(lineChange.Change, suppressed) {
					err := fmt.Errorf("file '%s' should not contains '%s'", file, suppressed)
					validationErrors = *multierror.Append(&validationErrors, err)

					break
				}
			}
		}
	}

	return validationErrors.ErrorOrNil()
}

func (rule *SuppressedText) GetPosition() byte {
	return PostScripts
}

func (rule *SuppressedText) Compile(variables map[string]any) {
	rule.BaseRule.Compile(variables)
	utils.FillTemplatesArray(rule.Substrings, variables)
	utils.FillTemplatesArray(rule.ExcludedGlobs, variables)
}
