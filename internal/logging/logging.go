package logging

import (
	"go.uber.org/zap"
)

func InitLogger() (*zap.Logger, error) {
	logger, err := zap.NewDevelopmentConfig().Build()
	return logger, err
}
