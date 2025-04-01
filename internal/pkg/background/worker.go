package background

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/iamsorryprincess/go-project-layout/internal/pkg/log"
)

type WorkerFunc func(ctx context.Context) error

type Worker struct {
	logger log.Logger

	mu sync.Mutex
	wg sync.WaitGroup
}

func NewWorker(logger log.Logger) *Worker {
	return &Worker{
		logger: logger,
	}
}

func (w *Worker) RunWithInterval(ctx context.Context, name string, interval time.Duration, worker WorkerFunc) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.wg.Add(1)

	go func(ctx context.Context, name string, interval time.Duration, worker WorkerFunc) {
		defer w.wg.Done()

		now := time.Now()
		if err := worker(ctx); err != nil {
			if errors.Is(err, context.Canceled) {
				w.logger.Info().Str("worker", name).Msg("worker cancelled")
				return
			}
			w.logger.Error().Err(err).Str("worker", name).Msg("worker failed")
		}
		w.logger.Debug().Str("worker", name).Str("time", time.Since(now).String()).Msg("worker done")

		timer := time.NewTimer(interval)
		defer timer.Stop()

		for {
			select {
			case <-timer.C:
				now = time.Now()
				if err := worker(ctx); err != nil {
					if errors.Is(err, context.Canceled) {
						w.logger.Info().Str("worker", name).Msg("worker cancelled")
						return
					}
					w.logger.Error().Err(err).Str("worker", name).Msg("worker failed")
				}
				w.logger.Debug().Str("worker", name).Str("time", time.Since(now).String()).Msg("worker done")
				timer.Reset(interval)
			case <-ctx.Done():
				w.logger.Info().Str("worker", name).Msg("worker cancelled")
				return
			}
		}
	}(ctx, name, interval, worker)
}

func (w *Worker) Wait() {
	w.wg.Wait()
}
