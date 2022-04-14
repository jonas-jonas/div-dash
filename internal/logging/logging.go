package logging

import (
	"fmt"

	"go.uber.org/zap"
)

type (
	LoggerProvider interface {
		Logger() *zap.SugaredLogger
	}
)

func NewLogger() *zap.SugaredLogger {
	newLogger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Printf("could not instantiate zap logger: %s", err.Error())
	}
	return newLogger.Sugar()
}
