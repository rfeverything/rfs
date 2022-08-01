package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	cfg *viper.Viper
}

func NewConfig() *Config {
	cfg := viper.New()
	cfg.SetConfigName("config")
	cfg.AddConfigPath(".")
	viper.AddConfigPath("/etc/rfs/")
	viper.AddConfigPath("$HOME/.rfs")
	cfg.SetConfigType("yaml")
	if err := cfg.ReadInConfig(); err != nil {
		panic(err)
	}

	return &Config{cfg: cfg}
}
