package handling

import (
	"context"
	"sync"

	"github.com/hashicorp/go-multierror"
)

type workersGroup struct {
	wg  sync.WaitGroup
	ctx context.Context
	mu  sync.Mutex
	err *multierror.Error
}

func workersGroupWithContext(ctx context.Context) *workersGroup {
	return &workersGroup{ctx: ctx}
}

func (g *workersGroup) Wait() error {
	g.wg.Wait()

	return g.err.ErrorOrNil()
}

func (g *workersGroup) Go(f func(context.Context) error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		if err := g.ctx.Err(); err != nil {
			g.error(err)

			return
		}

		if err := f(g.ctx); err != nil {
			g.error(err)
		}
	}()
}

func (g *workersGroup) error(err error) {
	g.mu.Lock()
	g.err = multierror.Append(g.err, err)
	g.mu.Unlock()
}
