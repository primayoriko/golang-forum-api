package logger

import (
	"os"
	"sync"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// Event is class as the representation of each requests that gonna be logged
type Event struct {
	id      int
	message string
}

// StandardLogger is class for logging that enforces specific formats message
type StandardLogger struct {
	*logrus.Logger
}

var (
	logger *StandardLogger
	once   sync.Once
)

// NewLogger is nethod for initialize the log functionality
func NewLogger() *StandardLogger {
	// logPath := os.Getenv("LOG_FILEPATH")
	pathMap := lfshook.PathMap{
		logrus.InfoLevel:  os.Getenv("INFO_LOG_FILEPATH"),
		logrus.ErrorLevel: os.Getenv("ERROR_LOG_FILEPATH"),
		logrus.PanicLevel: os.Getenv("PANIC_LOG_FILEPATH"),
		logrus.FatalLevel: os.Getenv("FATAL_LOG_FILEPATH"),
	}

	var baseLogger = logrus.New()
	baseLogger.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	))

	var standardLogger = &StandardLogger{baseLogger}

	logger = standardLogger

	return standardLogger
}

// GetInstance returning instance of logger
func GetInstance() *StandardLogger {
	once.Do(func() {
		logger = NewLogger()
	})

	return logger
}

// WriteLog used to writing log both to terminal and file
func (l *StandardLogger) WriteLog(id uint32, username string,
	method string, path string, statusCode int, statusText string) {
	fields := logrus.Fields{
		"id":          id,
		"username":    username,
		"method":      method,
		"path":        path,
		"status_code": statusCode,
	}

	if statusCode < 500 {
		l.WithFields(fields).Infoln(statusText)
	} else {
		l.WithFields(fields).Errorln(statusText)
	}
}
