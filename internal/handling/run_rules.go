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
	if h.WorkersCount == 0 {
		return errors.New("incorrect workers count")
	}

	if len(rules) == 0 {
		log.Debugf("no rules founded")

		return nil
	}

	input := make(chan configuration.Rule)
	output := startWorkers(ctx, input, utils.Min(h.WorkersCount, len(rules)))

	for _, rule := range rules {
		input <- rule
	}

	close(input)

	var multierr *multierror.Error

	for err := range output {
		multierr = multierror.Append(multierr, err)
	}

	return multierr.ErrorOrNil()
}

func startWorkers(ctx coxtext, input in, count int) chan error {
	wg := sync.WaitGroup{}

	output := make(chan error)

	wg.Add(count)

	for i := 0; i < count; i++ {
		go worker(i, &wg, ctx, input, output)
	}

	go func() {
		wg.Wait()
		log.Debug("all workers finished")
		close(output)
	}()

	return output
}

// TODO: Add panic interceptor.
func worker(id int, wg *sync.WaitGroup, ctx coxtext, input in, output out) {
	log.Debugf("workder %d started", id)
	defer log.Debugf("workder %d finished", id)
	defer wg.Done()

	for rule := range input {
		log.Debugf("workder %d received rules %s", id, rule.GetPrefix())
		err := checkRule(ctx, rule)
		// TODO: Move canclation to workers run method
		if err != nil {
			if !validation.IsValidationError(err) {
				log.Debugf("workder %d check rule with error %s", id, err)
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
