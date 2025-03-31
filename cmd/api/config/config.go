package config

import (
	"time"

	"github.com/iamsorryprincess/go-project-layout/internal/pkg/database/clickhouse"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/database/mysql"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/database/redis"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/http"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/messaging/nats"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/queue"
	"github.com/spf13/viper"
)

type Config struct {
	LogLevel string `mapstructure:"loglevel"`

	Mysql mysql.Config `mapstructure:"mysql"`

	Redis redis.Config `mapstructure:"redis"`

	Clickhouse clickhouse.Config `mapstructure:"clickhouse"`

	Nats nats.Config `mapstructure:"nats"`

	HTTP http.Config `mapstructure:"http"`

	TestQueue queue.MemoryBatchQueueConfig `mapstructure:"test_queue"`
}

func SetDefaults() {
	viper.SetDefault("loglevel", "info")

	viper.SetDefault("mysql.max_open_connections", 10)
	viper.SetDefault("mysql.max_idle_connections", 10)
	viper.SetDefault("mysql.connection_max_lifetime", time.Minute*5)
	viper.SetDefault("mysql.connection_max_idle_time", time.Minute*5)

	viper.SetDefault("http.port", 8080)

	viper.SetDefault("test_queue.name", "test")
	viper.SetDefault("test_queue.workers_count", 5)
	viper.SetDefault("test_queue.batch_size", 500)
	viper.SetDefault("test_queue.buffer_size", 500)
	viper.SetDefault("test_queue.handle_interval", "1s")
}
