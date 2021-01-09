package validators

import (
	"fisherman/internal"
	"fisherman/utils"
	"fmt"
	"regexp"
	"strings"
)

func MessageNotEmpty(ctx internal.SyncContext, notEmpty bool) error {
	if notEmpty && utils.IsEmpty(ctx.Message()) {
		return fmt.Errorf("commit message should not be empty")
	}

	return nil
}

func MessageHasPrefix(ctx internal.SyncContext, prefix string) error {
	if !utils.IsEmpty(prefix) && !strings.HasPrefix(ctx.Message(), prefix) {
		return fmt.Errorf("commit message should have prefix '%s'", prefix)
	}

	return nil
}

func MessageHasSuffix(ctx internal.SyncContext, suffix string) error {
	if !utils.IsEmpty(suffix) && !strings.HasSuffix(ctx.Message(), suffix) {
		return fmt.Errorf("commit message should have suffix '%s'", suffix)
	}

	return nil
}

func MessageRegexp(ctx internal.SyncContext, expression string) error {
	if !utils.IsEmpty(expression) {
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
