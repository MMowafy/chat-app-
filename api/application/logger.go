package application

import (
	"go.uber.org/zap"
	"log"
)

type Logger struct {
	logService *zap.SugaredLogger
}

func NewLogger() *Logger {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Printf("Failed to set new zap logger with err --> %s ", err.Error())
	}
	sLogger := zapLogger.Sugar()

	logger := &Logger{
		logService: sLogger,
	}
	return logger
}

// Debug uses fmt.Sprint to construct and log a message.
func (logger *Logger) Debug(args ...interface{}) {
	logger.logService.Debug(args)
}

// Info uses fmt.Sprint to construct and log a message.
func (logger *Logger) Info(args ...interface{}) {
	logger.logService.Info(args)
}

// Warn uses fmt.Sprint to construct and log a message.
func (logger *Logger) Warn(args ...interface{}) {
	logger.logService.Warn(args)
}

// Error uses fmt.Sprint to construct and log a message.
func (logger *Logger) Error(args ...interface{}) {
	logger.logService.Error(args)
}

// Error uses fmt.Sprint to construct and log a message.
func (logger *Logger) Panic(args ...interface{}) {
	logger.logService.Panic(args)
}

// Error uses fmt.Sprint to construct and log a message.
func (logger *Logger) Fatal(args ...interface{}) {
	logger.logService.Fatal(args)
}

// Debugf uses fmt.Sprintf to log a templated message.
func (logger *Logger) Debugf(template string, args ...interface{}) {
	logger.logService.Debugf(template, args)
}

// Infof uses fmt.Sprintf to log a templated message.
func (logger *Logger) Infof(template string, args ...interface{}) {
	logger.logService.Infof(template, args)
}

// Warnf uses fmt.Sprintf to log a templated message.
func (logger *Logger) Warnf(template string, args ...interface{}) {
	logger.logService.Warnf(template, args)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (logger *Logger) Errorf(template string, args ...interface{}) {
	logger.logService.Errorf(template, args)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (logger *Logger) Panicf(template string, args ...interface{}) {
	logger.logService.Panicf(template, args)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (logger *Logger) Fatalf(template string, args ...interface{}) {
	logger.logService.Fatalf(template, args)
}
