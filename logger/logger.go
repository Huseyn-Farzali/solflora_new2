package logger

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

var logger *log.Logger

func Init() {
	logger = log.New()
	logger.SetFormatter(&log.TextFormatter{})

	loggerLevel := os.Getenv("LOG_LEVEL")
	logger.SetLevel(parseLogLevel(loggerLevel))

	logger.Infof("[INFO] logger active with level: %s", loggerLevel)
}

func Logger() *log.Logger {
	return logger
}

func parseLogLevel(level string) log.Level {
	switch strings.ToLower(level) {
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "warn":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	case "panic":
		return log.PanicLevel
	default:
		return log.InfoLevel
	}
}
