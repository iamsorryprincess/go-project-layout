package mysql

import "time"

type Config struct {
	ConnectionString string `mapstructure:"connection_string"`

	MaxOpenConnections int `mapstructure:"max_open_connections"`
	MaxIdleConnections int `mapstructure:"max_idle_connections"`

	ConnectionMaxLifetime time.Duration `mapstructure:"connection_max_lifetime"`
	ConnectionMaxIdleTime time.Duration `mapstructure:"connection_max_idle_time"`
}
