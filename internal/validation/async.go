package validation

import (
	"context"
	"fisherman/infrastructure/log"
	"fmt"
	"sync"

	"github.com/hashicorp/go-multierror"
)

type AsyncValidationContext interface {
	SyncValidationContext
	context.Context
	Stop()
}

type AsyncValidator = func(ctx AsyncValidationContext) AsyncValidationResult

func RunAsync(ctx AsyncValidationContext, validators []AsyncValidator) error {
	output := make(chan AsyncValidationResult)

	var multierr *multierror.Error

	go runAsyncInternal(output, ctx, validators)

	for result := range output {
		if result.IsCanceled() {
			log.Infof("[%s] was skipped", result.Name)

			continue
		}

		if result.IsSuccessful() {
			log.Infof("[%s] complied (executed in %s)", result.Name, result.Time)
		} else {
			log.Infof("[%s] failed (executed in %s)", result.Name, result.Time)
			multierr = multierror.Append(multierr, fmt.Errorf("[%s] %s", result.Name, result.Error))
		}
	}

	return multierr.ErrorOrNil()
}

func runAsyncInternal(output chan AsyncValidationResult, ctx AsyncValidationContext, validators []AsyncValidator) {
	var wg sync.WaitGroup
	for _, validator := range validators {
		wg.Add(1)
		go wrap(output, &wg, validator)(ctx)
	}
	wg.Wait()
	close(output)
}

type wrappedValidator = func(ctx AsyncValidationContext)

func wrap(output chan AsyncValidationResult, wg *sync.WaitGroup, validator AsyncValidator) wrappedValidator {
	return func(ctx AsyncValidationContext) {
		defer wg.Done()
		result := validator(ctx)

		if context.Canceled == ctx.Err() {
			result.Error = ctx.Err()
		}

		if context.DeadlineExceeded == ctx.Err() {
			result.Error = ctx.Err()
		}

		output <- result

		if !result.IsSuccessful() {
			ctx.Stop()
		}
	}
}
