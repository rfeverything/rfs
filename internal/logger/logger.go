package logger

import (
	"github.com/rfeverything/rfs/internal/config"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func init() {
	switch config.Global().GetString("debug") {
	case "true":
		logger, _ = zap.NewDevelopment()
	default:
		logger, _ = zap.NewProduction()
	}
}

func Global() *zap.Logger {
	return logger
}
