// logger/logger.go
package logger

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

const (
	CtxLoggerKey = "logger"
)

var (
	log = logrus.New()
)

type Logger struct {
	entry *logrus.Entry
}

func init() {
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
}

func NewLogger(requestID string) *Logger {
	log.Info("checks")

	if requestID != "" {
		return &Logger{
			entry: log.WithField("request_id", requestID),
		}
	}
	return &Logger{entry: logrus.NewEntry(log)}
}

// FromCtx retrieves logger from context
func FromCtx(ctx context.Context) *Logger {
	if log, ok := ctx.Value(CtxLoggerKey).(*Logger); ok {
		return log
	}
	return NewLogger("")
}

func (l *Logger) Info(msg string, args ...interface{})  { l.entry.Info(fmt.Sprintf(msg, args...)) }
func (l *Logger) Error(msg string, args ...interface{}) { l.entry.Error(fmt.Sprintf(msg, args...)) }
