package config

import (
	"log"
	"sync"
)

var logger *log.Logger
var loggerOnce sync.Once

func initLogger() {
	logger = log.New(log.Writer(), "div-dash: ", log.Lshortfile|log.LstdFlags)
}

func Logger() *log.Logger {
	loggerOnce.Do(initLogger)
	return logger
}
