package log

import "go.uber.org/zap"

var (
	Logger *zap.Logger
)

func NewLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}
