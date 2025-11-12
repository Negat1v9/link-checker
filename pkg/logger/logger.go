package logger

import (
	"fmt"
	"log/slog"
	"os"
)

const (
	envProd  = "prod"
	envLocal = "local"
)

type Logger struct {
	l *slog.Logger
}

func NewLogger(env string) *Logger {
	l := newSlogLogger(env)
	return &Logger{
		l: l,
	}
}

func (l *Logger) Debugf(template string, args ...any) {
	l.l.Debug(fmt.Sprintf(template, args...))
}

func (l *Logger) Infof(template string, args ...any) {
	l.l.Info(fmt.Sprintf(template, args...))
}

func (l *Logger) Warnf(template string, args ...any) {
	l.l.Warn(fmt.Sprintf(template, args...))
}

func (l *Logger) Errorf(template string, args ...any) {
	l.l.Error(fmt.Sprintf(template, args...))
}

func newSlogLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {

	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		log.Info("loger info", slog.String("level", "[DEBUG]"))

	case envProd:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		log.Info("loger info", slog.String("level", "[INFO]"))
	default:
		panic("environment variable for logger not specified")
	}

	log.Info("logger", slog.String("environment", env))

	return log
}
