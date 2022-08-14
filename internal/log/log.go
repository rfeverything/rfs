package log

import (
	"github.com/rfeverything/rfs/internal/config"
	"go.uber.org/zap"
)

var (
	Logger *zap.Logger
)

func init() {
	switch config.Global().GetString("debug") {
	case "true":
		Logger, _ = zap.NewDevelopment()
	default:
		Logger, _ = zap.NewProduction()
	}
}

func Global() *zap.Logger {
	return Logger
}
