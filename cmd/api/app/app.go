package app

import (
	"context"

	"github.com/iamsorryprincess/go-project-layout/cmd/api/config"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/background"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/configutils"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/database/mysql"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/http"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/log"
)

const serviceName = "api"

type App struct {
	ctx context.Context

	logger log.Logger
	config config.Config

	mysqlConnection *mysql.Connection

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

	if a.mysqlConnection, err = mysql.NewConnection(a.logger, a.config.Mysql); err != nil {
		a.logger.Error().Err(err).Msg("failed to connect to mysql")
		return err
	}

	a.logger.Info().Msg("connected to mysql")

	return nil
}

func (a *App) initHTTP() {
	a.httpServer = http.NewServer(a.logger, a.config.HTTP, nil)
	a.httpServer.Start()
}

func (a *App) close() {
	if a.httpServer != nil {
		a.httpServer.Shutdown()
	}
	if a.mysqlConnection != nil {
		a.mysqlConnection.Close()
	}
}
