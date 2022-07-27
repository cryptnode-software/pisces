package lib

import (
	"go.uber.org/zap"
)

// Logger to log anything
type Logger interface {
	Info(msg string, keysAndValues ...interface{})

	Warn(msg string, keysAndValues ...interface{})

	Error(msg string, keysAndValues ...interface{})

	Debug(msg string, keysAndValues ...interface{})

	Fatal(msg string, keysAndValues ...interface{})

	With(args ...interface{}) Logger

	Close()
}

// Zapper struct to implement the logger interface
type Zapper struct {
	env string

	logger *zap.SugaredLogger
}

// NewZapper instantiates a new Zap logger
func NewZapper(env string) Logger {
	var logger *zap.Logger
	var err error

	switch env {
	case EnvDev:
		logger, err = zap.NewDevelopment()
	default:
		logger, err = zap.NewProduction()
	}

	if err != nil {
		panic(err)
	}

	// will figure out how to use the level and format later.
	return &Zapper{
		env:    env,
		logger: logger.Sugar(),
	}
}

// With for general information
func (z *Zapper) With(args ...interface{}) Logger {
	return &Zapper{
		env:    z.env,
		logger: z.logger.With(args...),
	}
}

// Info for general information
func (z *Zapper) Info(msg string, keysAndValues ...interface{}) {
	z.logger.Infow(msg, keysAndValues...)
}

// Warn for warnings
func (z *Zapper) Warn(msg string, keysAndValues ...interface{}) {
	z.logger.Warnw(msg, keysAndValues...)
}

// Error for errors
func (z *Zapper) Error(msg string, keysAndValues ...interface{}) {
	z.logger.Errorw(msg, keysAndValues...)
}

// Debug for any debugging related logs.
func (z *Zapper) Debug(msg string, keysAndValues ...interface{}) {
	if z.env == EnvProd {
		return
	}

	z.logger.Debugw(msg, keysAndValues...)
}

// Fatal for any fatal errors
func (z *Zapper) Fatal(msg string, keysAndValues ...interface{}) {
	z.logger.Fatalw(msg, keysAndValues...)
}

// Close handles any clean up that needs to be taken care of
func (z *Zapper) Close() {
	z.logger.Sync() // flushes buffer, if any
}
