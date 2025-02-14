package nats

import (
	"fmt"
	"strings"

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

		nats.ReconnectWait(nats.DefaultReconnectWait),
		nats.ReconnectJitter(nats.DefaultReconnectJitter, nats.DefaultReconnectJitterTLS),
		nats.Timeout(nats.DefaultTimeout),
		nats.PingInterval(nats.DefaultPingInterval),
		nats.MaxPingsOutstanding(nats.DefaultMaxPingOut),
		nats.SyncQueueLen(nats.DefaultMaxChanLen),
		nats.ReconnectBufSize(nats.DefaultReconnectBufSize),
		nats.DrainTimeout(nats.DefaultDrainTimeout),
		nats.FlusherTimeout(nats.DefaultFlusherTimeout),

		nats.ReconnectHandler(func(_ *nats.Conn) {
			logger.Info().Msg("reconnected to nats")
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

		nats.DiscoveredServersHandler(func(conn *nats.Conn) {
			logger.Info().
				Interface("known", conn.Servers()).
				Interface("discovered", conn.DiscoveredServers()).
				Msg("nats servers")
		}),
	}

	url := strings.Join(config.Servers, ",")

	conn, err := nats.Connect(url, options...)
	if err != nil {
		return nil, err
	}

	if !conn.IsConnected() {
		return nil, fmt.Errorf("could not establish connection to %s", url)
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
