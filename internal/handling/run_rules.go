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

type (
	in      = <-chan configuration.Rule
	out     = chan<- error
	coxtext = internal.ExecutionContext
)

func (h *HookHandler) runRules(ctx coxtext, rules []configuration.Rule) error {
	input := make(chan configuration.Rule)
	output, err := startWorkers(ctx, input, h.WorkersCount)
	if err != nil {
		return err
	}

	filteredRules := []configuration.Rule{}
	for _, rule := range rules {
		shouldAdd := true

		condition := rule.GetContition()
		if !utils.IsEmpty(condition) {
			shouldAdd, err = h.Engine.Eval(condition, h.GlobalVariables)
			if err != nil {
				close(input)

				return err
			}
		}

		if shouldAdd {
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

func startWorkers(ctx coxtext, input in, count int) (chan error, error) {
	wg := sync.WaitGroup{}

	if count <= 0 {
		return nil, errors.New("incorrect workers count")
	}

	output := make(chan error)

	wg.Add(count)
	// TODO: suppress spawn unused workers
	for i := 0; i < count; i++ {
		go worker(i, &wg, ctx, input, output)
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return output, nil
}

// TODO: Add panic interceptor.
func worker(id int, wg *sync.WaitGroup, ctx coxtext, input in, output out) {
	log.Debugf("workder %d started", id)
	defer log.Debugf("workder %d finished", id)
	defer wg.Done()

	for rule := range input {
		err := checkRule(ctx, rule)
		// TODO: Move canclation to workers run method
		if err != nil {
			if !validation.IsValidationError(err) {
				ctx.Cancel()
			}

			typeName := rule.GetType()
			output <- errors.Errorf("[%s] %s", typeName, err)
		}
	}
}

// TODO: Add more detailed validation result.
func checkRule(ctx coxtext, rule configuration.Rule) error {
	writer := ctx.Output()
	defer writer.Close()

	prefix := fmt.Sprintf("%s | ", rule.GetPrefix())
	prefixedWriter := prefixwriter.NewWriter(writer, prefix)

	return rule.Check(ctx, prefixedWriter)
}
