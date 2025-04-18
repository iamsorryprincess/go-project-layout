package app

import (
	"context"
	"time"

	"github.com/iamsorryprincess/go-project-layout/cmd/api/config"
	httpapp "github.com/iamsorryprincess/go-project-layout/cmd/api/http"
	"github.com/iamsorryprincess/go-project-layout/cmd/api/model"
	"github.com/iamsorryprincess/go-project-layout/cmd/api/service"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/background"
	memorycache "github.com/iamsorryprincess/go-project-layout/internal/pkg/cache/memory"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/configutils"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/database/clickhouse"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/database/mysql"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/database/redis"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/http"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/log"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/messaging/nats"
	memoryqueue "github.com/iamsorryprincess/go-project-layout/internal/pkg/queue/memory"
	redisqueue "github.com/iamsorryprincess/go-project-layout/internal/pkg/queue/redis"
)

const serviceName = "api"

type App struct {
	ctx context.Context

	logger log.Logger
	config config.Config

	mysqlConn      *mysql.Connection
	redisConn      *redis.Connection
	clickhouseConn *clickhouse.Connection
	natsConn       *nats.Connection

	testCache *memorycache.Cache[string, int]

	testQueue *memoryqueue.Queue[int]

	userService *service.UserService

	testConsumer *redisqueue.Consumer[model.User]

	worker *background.Worker

	httpServer *http.Server
}

func New() *App {
	return &App{}
}

func (a *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())

	defer a.close()
	defer cancel()

	a.ctx = ctx

	if err := a.initConfig(); err != nil {
		return
	}

	if err := a.initDatabases(); err != nil {
		return
	}

	if err := a.initNats(); err != nil {
		return
	}

	a.initServices()

	a.initHTTP()

	a.logger.Info().Msg("service started")

	s := background.Wait()

	a.logger.Info().Str("stop_signal", s.String()).Msg("service stopped")
}

func (a *App) initConfig() error {
	var err error
	a.config, err = configutils.Parse[config.Config](config.SetDefaults)
	if err != nil {
		log.New("error", serviceName).Error().Err(err).Msg("failed to parse config")
		return err
	}
	a.logger = log.New(a.config.LogLevel, serviceName)
	return nil
}

func (a *App) initDatabases() error {
	var err error

	if a.mysqlConn, err = mysql.NewConnection(a.logger, a.config.Mysql); err != nil {
		a.logger.Error().Err(err).Msg("failed to connect to mysql")
		return err
	}

	a.logger.Info().Msg("connected to mysql")

	if a.redisConn, err = redis.NewConnection(a.logger, a.config.Redis); err != nil {
		a.logger.Error().Err(err).Msg("failed to connect to redis")
		return err
	}

	a.logger.Info().Msg("connected to redis")

	if a.clickhouseConn, err = clickhouse.NewConnection(a.logger, a.config.Clickhouse); err != nil {
		a.logger.Error().Err(err).Msg("failed to connect to clickhouse")
		return err
	}

	a.logger.Info().Msg("connected to clickhouse")

	return nil
}

func (a *App) initNats() error {
	a.config.Nats.Name = serviceName
	var err error
	if a.natsConn, err = nats.NewConnection(a.logger, a.config.Nats); err != nil {
		a.logger.Error().Err(err).Msg("failed to connect to nats")
		return err
	}

	a.logger.Info().Msg("connected to nats")

	return nil
}

func (a *App) initServices() {
	a.testCache = memorycache.NewCache[string, int](a.ctx, memorycache.Config{
		ClearInterval: time.Second,
	})

	_ = a.testCache.Set(a.ctx, "test", time.Millisecond*300, 1)

	a.testQueue = memoryqueue.NewQueue[int](a.ctx, a.logger, a.config.TestQueue, nil)

	a.userService = service.NewUserService(a.logger)

	a.testConsumer = redisqueue.NewConsumer(a.logger, "test", 10, a.redisConn, a.userService)

	a.worker = background.NewWorker(a.logger)
	a.worker.RunWithInterval(a.ctx, "test", time.Second, a.testConsumer.Consume)
}

func (a *App) initHTTP() {
	handler := httpapp.NewHandler(a.logger)
	a.httpServer = http.NewServer(a.logger, a.config.HTTP, handler)
	a.httpServer.Start()
}

func (a *App) close() {
	if a.httpServer != nil {
		a.httpServer.Shutdown()
	}
	if a.worker != nil {
		a.worker.Wait()
	}
	if a.testQueue != nil {
		a.testQueue.Close()
	}
	if a.testCache != nil {
		a.testCache.Wait()
	}
	if a.natsConn != nil {
		a.natsConn.Shutdown()
	}
	if a.mysqlConn != nil {
		a.mysqlConn.Close()
	}
	if a.redisConn != nil {
		a.redisConn.Close()
	}
	if a.clickhouseConn != nil {
		a.clickhouseConn.Close()
	}
}
