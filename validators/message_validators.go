package validators

import (
	"fisherman/internal/validation"
	"fisherman/utils"
	"fmt"
	"regexp"
	"strings"
)

func MessageNotEmpty(ctx validation.SyncValidationContext, notEmpty bool) error {
	if notEmpty && utils.IsEmpty(ctx.Message()) {
		return fmt.Errorf("commit message should not be empty")
	}

	return nil
}

func MessageHasPrefix(ctx validation.SyncValidationContext, prefix string) error {
	if utils.IsNotEmpty(prefix) && !strings.HasPrefix(ctx.Message(), prefix) {
		return fmt.Errorf("commit message should have prefix '%s'", prefix)
	}

	return nil
}

func MessageHasSuffix(ctx validation.SyncValidationContext, suffix string) error {
	if utils.IsNotEmpty(suffix) && !strings.HasSuffix(ctx.Message(), suffix) {
		return fmt.Errorf("commit message should have suffix '%s'", suffix)
	}

	return nil
}

func MessageRegexp(ctx validation.SyncValidationContext, expression string) error {
	if utils.IsNotEmpty(expression) {
		matched, err := regexp.MatchString(expression, ctx.Message())
		if err != nil {
			return err
		}

		if !matched {
			return fmt.Errorf("commit message should be matched regular expression '%s'", expression)
		}
	}

	return nil
}
