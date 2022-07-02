package executors

import (
	"context"
	"github.com/Lameaux/core/logger"
	"github.com/Lameaux/core/runners"
	"github.com/Lameaux/core/workers"
	"golang.org/x/sync/errgroup"
)

type Executor struct {
	ErrGroup      *errgroup.Group
	CancelWorkers context.CancelFunc
}

func NewExecutor(wrks []workers.Worker) *Executor {
	ctx, cancel := context.WithCancel(context.Background())

	rnrs := make([]*runners.Runner, 0)

	for _, w := range wrks {
		rnrs = append(rnrs, runners.NewRunner(ctx, w))
	}

	var wg errgroup.Group

	for _, r := range rnrs {
		r := r

		wg.Go(func() error {
			return r.Exec()
		})
	}

	return &Executor{&wg, cancel}
}

func (e *Executor) Shutdown() {
	logger.Infow("shutting down workers")

	e.CancelWorkers()

	if err := e.ErrGroup.Wait(); err != nil {
		logger.Errorw("error while stopping workers", "error", err)
	}

	logger.Infow("workers stopped")
}
