package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"strings"
	"time"
)

// Logger -.
type Logger interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

// ZeroLogLogger -.
type ZeroLogLogger struct {
	logger *zerolog.Logger
}

// New -.
func NewZeroLog(level string) Logger {
	var l zerolog.Level

	switch strings.ToLower(level) {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)

	skipFrameCount := 3
	var logger zerolog.Logger
	if l == zerolog.DebugLevel {
		logger = zerolog.New(os.Stdout).Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
			With().
			Timestamp().
			CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).
			Logger()
	} else {
		logger = zerolog.New(os.Stdout).
			With().
			Timestamp().
			CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).
			Logger()
	}

	return &ZeroLogLogger{
		logger: &logger,
	}
}

func (l *ZeroLogLogger) formatMessage(message any) string {
	switch t := message.(type) {
	case error:
		return t.Error()
	case string:
		return t
	default:
		return fmt.Sprintf("Unknown type %v", message)
	}
}

// Debug -.
func (l *ZeroLogLogger) Debug(message any, args ...any) {
	mf := l.formatMessage(message)
	l.log(l.logger.Debug(), mf, args...)
}

// Info -.
func (l *ZeroLogLogger) Info(message string, args ...any) {
	mf := l.formatMessage(message)
	l.log(l.logger.Info(), mf, args...)
}

// Warn -.
func (l *ZeroLogLogger) Warn(message string, args ...any) {
	mf := l.formatMessage(message)
	l.log(l.logger.Warn(), mf, args...)
}

// Error -.
func (l *ZeroLogLogger) Error(message interface{}, args ...any) {
	mf := l.formatMessage(message)
	l.log(l.logger.Error(), mf, args...)
}

// Fatal -.
func (l *ZeroLogLogger) Fatal(message interface{}, args ...any) {
	mf := l.formatMessage(message)
	l.log(l.logger.Fatal(), mf, args...)

	os.Exit(1)
}

func (l *ZeroLogLogger) log(e *zerolog.Event, m string, args ...any) {
	if len(args) == 0 {
		e.Msg(m)
	} else {
		e.Msgf(m, args)
	}
}
