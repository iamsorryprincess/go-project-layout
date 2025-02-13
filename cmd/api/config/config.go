package config

import (
	"time"

	"github.com/iamsorryprincess/go-project-layout/internal/pkg/database/clickhouse"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/database/mysql"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/database/redis"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/http"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/messaging/nats"
	"github.com/spf13/viper"
)

type Config struct {
	LogLevel string `mapstructure:"loglevel"`

	Mysql mysql.Config `mapstructure:"mysql"`

	Redis redis.Config `mapstructure:"redis"`

	Clickhouse clickhouse.Config `mapstructure:"clickhouse"`

	Nats nats.Config `mapstructure:"nats"`

	HTTP http.Config `mapstructure:"http"`
}

func SetDefaults() {
	viper.SetDefault("loglevel", "info")

	viper.SetDefault("mysql.max_open_connections", 10)
	viper.SetDefault("mysql.max_idle_connections", 10)
	viper.SetDefault("mysql.connection_max_lifetime", time.Minute*5)
	viper.SetDefault("mysql.connection_max_idle_time", time.Minute*5)

	viper.SetDefault("http.port", 8080)
}
