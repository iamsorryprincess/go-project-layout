package clickhouse

import (
	"context"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/log"
)

type Connection struct {
	logger log.Logger
	config Config
	driver.Conn
}

func NewConnection(logger log.Logger, config Config) (*Connection, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: config.Hosts,
		Auth: clickhouse.Auth{
			Database: config.Database,
			Username: config.User,
			Password: config.Password,
		},

		Debug: config.Debug,
		Debugf: func(format string, v ...any) {
			logger.Debug().Msgf(format, v...)
		},

		Settings: clickhouse.Settings{
			"max_execution_time": config.MaxExecutionTime,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:     config.DialTimeout,
		MaxOpenConns:    config.MaxOpenConnection,
		MaxIdleConns:    config.MaxIdleConnection,
		ConnMaxLifetime: config.MaxLifeConnection,

		ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
	})
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	return &Connection{
		logger: logger,
		config: config,
		Conn:   conn,
	}, nil
}

func (c *Connection) Close() {
	if err := c.Conn.Close(); err != nil {
		c.logger.Error().Err(err).Msg("failed to close clickhouse connection")
	}
}

func (c *Connection) CloseRows(rows driver.Rows) {
	if err := rows.Close(); err != nil {
		c.logger.Error().Err(err).Msg("failed to close clickhouse rows")
	}
}
