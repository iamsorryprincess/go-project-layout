package memory

import "time"

type Config struct {
	Name         string `mapstructure:"name"`
	WorkersCount int    `mapstructure:"workers_count"`

	BatchSize  int `mapstructure:"batch_size"`
	BufferSize int `mapstructure:"buffer_size"`

	HandleInterval time.Duration `mapstructure:"handle_interval"`
}
