package handling

import (
	"fisherman/internal"
	"fisherman/internal/configuration"
	"fisherman/internal/utils"
	"fisherman/pkg/log"
	"fisherman/pkg/prefixwriter"
	"fmt"

	"github.com/go-errors/errors"
)

type (
	in      = <-chan configuration.Rule
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

	input := rulesReduser(ctx, rules)

	g := workersGroupWithContext(ctx)
	for i := 0; i < utils.Min(h.WorkersCount, len(rules)); i++ {
		g.Go(worker(i, input))
	}

	return g.Wait()
}

func rulesReduser(ctx coxtext, rules []configuration.Rule) chan configuration.Rule {
	input := make(chan configuration.Rule)

	go func() {
		defer close(input)
		for _, rule := range rules {
			if err := ctx.Err(); err != nil {
				return
			}

			input <- rule
		}
	}()

	return input
}

// TODO: Add panic interceptor.
func worker(id int, input in) func(coxtext) error {
	return func(ctx coxtext) error {
		log.Debugf("workder %d started", id)
		defer log.Debugf("workder %d finished", id)

		for rule := range input {
			log.Debugf("workder %d received rules %s", id, rule.GetPrefix())
			err := checkRule(ctx, rule)
			// TODO: Move canclation to workers run method
			if err != nil {
				return errors.Errorf("[%s] %s", rule.GetType(), err)
			}
		}

		return nil
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
