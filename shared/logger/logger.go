package logger

import (
	"github.com/adcmdev/logger"
)

func Init(logLevel string) {
	level := logger.LevelFromString(logLevel)

	logger.New(level)
}

func CurrentLevel() string {
	level := logger.CurrentLevel()
	return logger.LevelToString(level)
}

func SetLevel(level string) {
	logger.SetLevel(logger.LevelFromString(level))
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}
