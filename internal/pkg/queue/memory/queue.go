package memory

import (
	"context"
	"sync"
	"time"

	"github.com/iamsorryprincess/go-project-layout/internal/pkg/log"
)

type BatchHandler[T any] interface {
	Handle(ctx context.Context, batch []T) error
}

type Queue[T any] struct {
	logger log.Logger

	name string

	ch chan T
	wg sync.WaitGroup

	handler BatchHandler[T]
}

func NewQueue[T any](ctx context.Context, logger log.Logger, config Config, handler BatchHandler[T]) *Queue[T] {
	if config.BufferSize <= 0 {
		config.BufferSize = 1
	}

	queue := &Queue[T]{
		logger:  logger,
		name:    config.Name,
		ch:      make(chan T, config.BufferSize),
		handler: handler,
	}

	for i := 0; i < config.WorkersCount; i++ {
		queue.wg.Add(1)
		go func(ctx context.Context, workerID int, queue *Queue[T]) {
			defer queue.wg.Done()

			timer := time.NewTimer(config.HandleInterval)
			batch := make([]T, 0, config.BatchSize)

			for {
				select {
				case message, ok := <-queue.ch:
					if !ok {
						timer.Stop()
						queue.drainMessages(ctx, workerID, batch)
						return
					}
					batch = append(batch, message)
					if len(batch) >= config.BatchSize {
						queue.processMessages(ctx, workerID, batch)
						batch = batch[:0]
						timer.Reset(config.HandleInterval)
					}
				case <-timer.C:
					if len(batch) > 0 {
						queue.processMessages(ctx, workerID, batch)
						batch = batch[:0]
					}
					timer.Reset(config.HandleInterval)
				case <-ctx.Done():
					timer.Stop()
					queue.drainMessages(ctx, workerID, batch)
					return
				}
			}
		}(ctx, i, queue)
	}

	return queue
}

func (q *Queue[T]) Push(message T) error {
	q.ch <- message
	return nil
}

func (q *Queue[T]) Close() {
	close(q.ch)
	q.wg.Wait()
	q.logger.Info().Str("queue", q.name).Msg("memory batch queue closed")
}

func (q *Queue[T]) processMessages(ctx context.Context, workerID int, messages []T) {
	if err := q.handler.Handle(ctx, messages); err != nil {
		q.logger.Error().
			Str("queue", q.name).
			Int("worker", workerID).
			Err(err).
			Msg("handle batch failed")
	}
}

func (q *Queue[T]) drainMessages(ctx context.Context, workerID int, messages []T) {
	if len(messages) > 0 {
		q.processMessages(ctx, workerID, messages)
	}

	drainedMessages := make([]T, 0, len(q.ch))
	for message := range q.ch {
		drainedMessages = append(drainedMessages, message)
	}

	if len(drainedMessages) > 0 {
		q.processMessages(ctx, workerID, drainedMessages)
	}
}
