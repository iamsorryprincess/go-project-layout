package queue

import (
	"context"
	"sync"
	"time"

	"github.com/iamsorryprincess/go-project-layout/internal/pkg/log"
)

type MemoryBatchQueueConfig struct {
	Name         string `mapstructure:"name"`
	WorkersCount int    `mapstructure:"workers_count"`

	BatchSize  int `mapstructure:"batch_size"`
	BufferSize int `mapstructure:"buffer_size"`

	HandleInterval time.Duration `mapstructure:"handle_interval"`
}

type BatchHandler[T any] interface {
	Handle(ctx context.Context, batch []T) error
}

type MemoryBatchQueue[T any] struct {
	ctx    context.Context
	logger log.Logger

	name string

	ch chan T
	wg sync.WaitGroup

	handler BatchHandler[T]
}

func NewMemoryBatchQueue[T any](ctx context.Context, logger log.Logger, config MemoryBatchQueueConfig, handler BatchHandler[T]) *MemoryBatchQueue[T] {
	if config.BufferSize <= 0 {
		config.BufferSize = 1
	}

	queue := &MemoryBatchQueue[T]{
		ctx:     context.Background(),
		logger:  logger,
		name:    config.Name,
		ch:      make(chan T, config.BufferSize),
		handler: handler,
	}

	for i := 0; i < config.WorkersCount; i++ {
		queue.wg.Add(1)
		go func(ctx context.Context, workerID int, queue *MemoryBatchQueue[T]) {
			defer queue.wg.Done()

			timer := time.NewTimer(config.HandleInterval)
			batch := make([]T, 0, config.BatchSize)

			for {
				select {
				case message, ok := <-queue.ch:
					if !ok {
						timer.Stop()
						queue.drainMessages(workerID, batch)
						return
					}
					batch = append(batch, message)
					if len(batch) >= config.BatchSize {
						queue.processMessages(workerID, batch)
						batch = batch[:0]
						timer.Reset(config.HandleInterval)
					}
				case <-timer.C:
					if len(batch) > 0 {
						queue.processMessages(workerID, batch)
						batch = batch[:0]
					}
					timer.Reset(config.HandleInterval)
				case <-ctx.Done():
					timer.Stop()
					queue.drainMessages(workerID, batch)
					return
				}
			}
		}(ctx, i, queue)
	}

	return queue
}

func (q *MemoryBatchQueue[T]) Push(message T) error {
	q.ch <- message
	return nil
}

func (q *MemoryBatchQueue[T]) Stop() {
	close(q.ch)
	q.wg.Wait()
	q.logger.Info().Str("queue", q.name).Msg("memory batch queue stopped")
}

func (q *MemoryBatchQueue[T]) processMessages(workerID int, batch []T) {
	if err := q.handler.Handle(q.ctx, batch); err != nil {
		q.logger.Error().
			Str("queue", q.name).
			Int("worker", workerID).
			Err(err).
			Msg("handle batch failed")
	}
}

func (q *MemoryBatchQueue[T]) drainMessages(workerID int, batch []T) {
	if len(batch) > 0 {
		q.processMessages(workerID, batch)
	}

	drainedMessages := make([]T, 0, len(q.ch))
	for message := range q.ch {
		drainedMessages = append(drainedMessages, message)
	}

	if len(drainedMessages) > 0 {
		q.processMessages(workerID, drainedMessages)
	}
}
