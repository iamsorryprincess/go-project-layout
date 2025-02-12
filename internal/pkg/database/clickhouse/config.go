package clickhouse

import "time"

type Config struct {
	Hosts    []string `mapstructure:"hosts"`
	Database string   `mapstructure:"database"`
	User     string   `mapstructure:"user"`
	Password string   `mapstructure:"password"`

	Debug       bool          `mapstructure:"debug"`
	DialTimeout time.Duration `mapstructure:"dial_timeout"`

	MaxOpenConnection int           `mapstructure:"max_open_connection"`
	MaxIdleConnection int           `mapstructure:"max_idle_connection"`
	MaxExecutionTime  int           `mapstructure:"max_execution_time"`
	MaxLifeConnection time.Duration `mapstructure:"max_life_connection"`
}
