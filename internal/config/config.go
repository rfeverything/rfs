package config

import (
	"github.com/spf13/viper"
)

var cfg *viper.Viper

func init() {
	cfg = viper.New()
	cfg.SetConfigName("config")
	cfg.AddConfigPath(".")
	viper.AddConfigPath("/etc/rfs/")
	viper.AddConfigPath("$HOME/.rfs")
	cfg.SetConfigType("yaml")
	if err := cfg.ReadInConfig(); err != nil {
		panic(err)
	}
}

func Global() *viper.Viper {
	return cfg
}
