package config

import (
	"github.com/spf13/viper"
)

var cfg *viper.Viper

func init() {
	cfg = viper.New()
	cfg.SetConfigName("config")
	cfg.AddConfigPath(".")
	cfg.AddConfigPath("/etc/rfs/")
	cfg.AddConfigPath("$HOME/.rfs")
	cfg.SetConfigType("yaml")

	cfg.SetDefault("client", map[string]interface{}{
		"server": "localhost",
		"port":   "8080",
	})
	cfg.SetDefault("metaserver", map[string]interface{}{
		"port": "8080",
		"host": "",
	})
	cfg.SetDefault("volume", map[string]interface{}{
		"port":   "8081",
		"host":   "",
		"dbpath": "./db",
	})
	cfg.SetDefault("metrics", map[string]interface{}{
		"port": "8082",
		"host": "",
	})

	if err := cfg.ReadInConfig(); err != nil {
		panic(err)
	}

}

func Global() *viper.Viper {
	return cfg
}
