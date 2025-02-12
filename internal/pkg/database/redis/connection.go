package redis

import (
	"context"
	"fmt"

	"github.com/iamsorryprincess/go-project-layout/internal/pkg/log"
	"github.com/redis/go-redis/v9"
)

type Connection struct {
	logger log.Logger
	config Config
	*redis.Client
}

func NewConnection(logger log.Logger, config Config) (*Connection, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Username: config.User,
		Password: config.Password,
		DB:       config.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &Connection{
		logger: logger,
		config: config,
		Client: client,
	}, nil
}

func (c *Connection) Close() {
	if err := c.Client.Close(); err != nil {
		c.logger.Error().Err(err).Msg("failed to close redis connection")
	}
}
