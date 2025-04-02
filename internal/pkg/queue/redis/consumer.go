package redis

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/iamsorryprincess/go-project-layout/internal/pkg/database/redis"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/log"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/queue"
)

type Consumer[T any] struct {
	logger log.Logger

	key   string
	count int

	conn    *redis.Connection
	handler queue.BatchHandler[T]
}

func NewConsumer[T any](logger log.Logger, key string, count int, conn *redis.Connection, handler queue.BatchHandler[T]) *Consumer[T] {
	return &Consumer[T]{
		logger:  logger,
		key:     key,
		count:   count,
		conn:    conn,
		handler: handler,
	}
}

func (c *Consumer[T]) Consume(ctx context.Context) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		result, err := c.conn.LPopCount(ctx, c.key, c.count).Result()
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return err
			}
			if !errors.Is(err, redis.ErrNil) {
				c.logger.Error().Str("key", c.key).Err(err).Msg("redis consumer LPOP failed")
			}
			return nil
		}

		if len(result) == 0 {
			return nil
		}

		data := make([]T, len(result))
		for i, item := range result {
			if err = json.Unmarshal([]byte(item), &data[i]); err != nil {
				return fmt.Errorf("redis consumer json unmarshal failed: %w", err)
			}
		}

		if err = c.handler.Handle(ctx, data); err != nil {
			if errors.Is(err, context.Canceled) {
				returnData := make([]interface{}, len(data))
				for i, item := range result {
					returnData[i] = item
				}
				if rErr := c.conn.RPush(context.Background(), c.key, returnData...).Err(); rErr != nil {
					c.logger.Error().Str("key", c.key).Err(rErr).Msg("redis consumer RPUSH data back failed")
				}
				return err
			}
			return fmt.Errorf("redis consumer handle failed: %w", err)
		}
	}
}

func parse1[T any](items []string) ([]T, error) {
	result := make([]T, len(items))
	for i, item := range items {
		if err := json.Unmarshal([]byte(item), &result[i]); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func parse2[T any](items []string) ([]T, error) {
	var buf bytes.Buffer
	buf.WriteString("[")
	n := len(items)
	for i := 0; i < n-1; i++ {
		buf.WriteString(items[i])
		buf.WriteString(",")
	}

	buf.WriteString(items[n-1])
	buf.WriteString("]")

	result := make([]T, 0, len(items))
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		return nil, err
	}

	return result, nil
}
