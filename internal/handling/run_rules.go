package handling

import (
	"errors"
	"fisherman/configuration"
	"fisherman/infrastructure/log"
	"fisherman/internal"
	"fisherman/internal/prefixwriter"
	"fmt"
	"os"
	"sync"

	"github.com/hashicorp/go-multierror"
)

type in = <-chan configuration.Rule
type out = chan<- error
type coxtext = internal.ExecutionContext

func (handler *HookHandler) runRules(ctx coxtext, rules []configuration.Rule) error {
	input := make(chan configuration.Rule)
	output := make(chan error)

	err := startWorkers(ctx, input, output, handler.WorkersCount)
	if err != nil {
		return err
	}

	filteredRules := []configuration.Rule{}
	for _, rule := range rules {
		condition, err := handler.Engine.Eval(rule.GetContition())
		if err != nil {
			close(input)

			return err
		}

		if condition {
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

	for i := 0; i < count; i++ {
		wg.Add(1)
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
		prefix := fmt.Sprintf("[%s]", rule.GetType())
		writer := prefixwriter.New(os.Stdout, prefix)
		err := rule.Check(writer, ctx)
		if err != nil {
			typeName := rule.GetType()
			output <- fmt.Errorf("[%s] %s", typeName, err)
		}
	}
}