package handling

import (
	"fisherman/internal"
	"fisherman/internal/configuration"
	"fisherman/internal/utils"
	"fisherman/internal/validation"
	"fisherman/pkg/log"
	"fisherman/pkg/prefixwriter"
	"fmt"
	"sync"

	"github.com/go-errors/errors"

	"github.com/hashicorp/go-multierror"
)

type in = <-chan configuration.Rule
type out = chan<- error
type coxtext = internal.ExecutionContext

func (h *HookHandler) runRules(ctx coxtext, rules []configuration.Rule) error {
	input := make(chan configuration.Rule)
	output := make(chan error)

	err := startWorkers(ctx, input, output, h.WorkersCount)
	if err != nil {
		return err
	}

	filteredRules := []configuration.Rule{}
	for _, rule := range rules {
		shouldAadd := true

		condition := rule.GetContition()
		if !utils.IsEmpty(condition) {
			shouldAadd, err = h.Engine.Eval(condition, h.GlobalVariables)
			if err != nil {
				close(input)

				return err
			}
		}

		if shouldAadd {
			filteredRules = append(filteredRules, rule)
		}
	}

	for _, rule := range filteredRules {
		input <- rule
	}

	close(input)

	var multierr *multierror.Error

	for err := range output {
		multierr = multierror.Append(multierr, err)
	}

	return multierr.ErrorOrNil()
}

func startWorkers(ctx coxtext, input in, output out, count int) error {
	wg := sync.WaitGroup{}

	if count <= 0 {
		return errors.New("incorrect workers count")
	}

	wg.Add(count)
	for i := 0; i < count; i++ {
		go worker(i, &wg, ctx, input, output)
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return nil
}

func worker(id int, wg *sync.WaitGroup, ctx coxtext, input in, output out) {
	log.Debugf("workder %d started", id)
	defer log.Debugf("workder %d finished", id)
	defer wg.Done()

	for rule := range input {
		prefix := fmt.Sprintf("%s | ", rule.GetPrefix())
		writer := prefixwriter.New(ctx.Output(), prefix)
		err := rule.Check(ctx, writer)
		if err != nil {
			if !validation.IsValidationError(err) {
				ctx.Cancel()
			}

			typeName := rule.GetType()
			output <- errors.Errorf("[%s] %s", typeName, err)
		}
	}
}
