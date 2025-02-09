package configutils

import (
	"flag"

	"github.com/spf13/viper"
)

type SetDefaultsFunc func()

func Parse[TConfig any](defaultsFunc SetDefaultsFunc) (TConfig, error) {
	path := flag.String("c", "config.json", "path to config file")
	flag.Parse()

	viper.SetConfigFile(*path)
	viper.SetConfigType("json")

	var cfg TConfig

	if err := viper.ReadInConfig(); err != nil {
		return cfg, err
	}

	if defaultsFunc != nil {
		defaultsFunc()
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
