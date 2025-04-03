package memory

import "time"

type Config struct {
	TTL           time.Duration `mapstructure:"ttl"`
	ClearInterval time.Duration `mapstructure:"clear_interval"`
}
