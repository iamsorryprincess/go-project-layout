package configutils

import (
	"flag"
	"strings"

	"github.com/spf13/viper"
)

type SetDefaultsFunc func()

func Parse[TConfig any](defaultsFunc SetDefaultsFunc, path ...string) (TConfig, error) {
	configPath := flag.String("c", "config.json", "path to config file")
	flag.Parse()

	if len(path) > 0 {
		p := strings.Join(path, "/") + "/" + *configPath
		configPath = &p
	}

	viper.SetConfigFile(*configPath)
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
