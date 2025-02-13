package nats

import (
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/log"
	"github.com/nats-io/nats.go"
)

type Connection struct {
	logger log.Logger
	*nats.Conn
}

func NewConnection(logger log.Logger, config Config) (*Connection, error) {
	options := []nats.Option{
		nats.Name(config.Name),

		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),

		nats.ReconnectHandler(func(_ *nats.Conn) {
			logger.Info().Msg("reconnected to nats")
		}),

		nats.ConnectHandler(func(_ *nats.Conn) {
			logger.Info().Msg("connected to nats")
		}),

		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			logger.Error().Err(err).Msg("nats disconnected")
		}),

		nats.ClosedHandler(func(_ *nats.Conn) {
			logger.Error().Msg("nats connection closed")
		}),

		nats.ErrorHandler(func(_ *nats.Conn, subscription *nats.Subscription, err error) {
			logger.Error().Err(err).Str("nats_subject", subscription.Subject).Msg("nats error")
		}),
	}

	conn, err := nats.Connect(config.URL, options...)
	if err != nil {
		return nil, err
	}

	return &Connection{
		logger: logger,
		Conn:   conn,
	}, nil
}

func (c *Connection) Shutdown() {
	if err := c.Conn.Drain(); err != nil {
		c.logger.Error().Err(err).Msg("nats drain error")
	}
}
