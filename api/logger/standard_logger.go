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

	invalidArgMessage      = Event{1, "Invalid arg: %s"}
	invalidArgValueMessage = Event{2, "Invalid value for argument: %s: %v"}
	missingArgMessage      = Event{3, "Missing arg: %s"}
	successMessage         = Event{4, "Missing arg: %s"}
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

	// log.SetFormatter(&log.JSONFormatter{})
	// standardFields := log.Fields{
	// 	"hostname": "staging-1",
	// 	"appname":  "foo-app",
	// 	"session":  "1ce3f6v",
	// }

	// standardLogger.WithFields(standardFields)
	// 	WithFields(
	// 		log.Fields{
	// 			"string": "foo",
	// 			"int":    1,
	// 			"float":  1.1,
	// 		}).
	// 	Info("Testing log with golang")

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

// WriteLog x
func WriteLog() {

}

// InvalidArg is a standard error message
func (l *StandardLogger) InvalidArg(argumentName string) {
	l.Errorf(invalidArgMessage.message, argumentName)
}

// InvalidArgValue is a standard error message
func (l *StandardLogger) InvalidArgValue(argumentName string, argumentValue string) {
	l.Errorf(invalidArgValueMessage.message, argumentName, argumentValue)
}

// MissingArg is a standard error message
func (l *StandardLogger) MissingArg(argumentName string) {
	l.Errorf(missingArgMessage.message, argumentName)
}

// SuccessArg is a standard error message
func (l *StandardLogger) SuccessArg(argumentName string) {
	l.Infof(missingArgMessage.message, argumentName)
}
