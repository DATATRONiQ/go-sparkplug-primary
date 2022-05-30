package util

import (
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func InitLogger(logFormat string, logFile string, logLevel string) {
	logrus.SetFormatter(&logrus.TextFormatter{})

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		panic(err)
	}
	logrus.SetLevel(level)
	switch logFormat {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})
	default:
		logrus.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: time.RFC3339,
		})
	}
	if logFile != "" {
		logrus.Info("Logging to file: ", logFile)
		logFile, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logrus.SetOutput(io.MultiWriter(logFile, os.Stdout))
		} else {
			logrus.Info("Failed to log to file, using default stderr")
		}
	}
}
