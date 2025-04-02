package queue

import "context"

type BatchHandler[T any] interface {
	Handle(ctx context.Context, batch []T) error
}
