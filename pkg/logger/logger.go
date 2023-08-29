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
	Info(message interface{}, args ...interface{})
	Warn(message interface{}, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

// ZeroLogLogger -.
type ZeroLogLogger struct {
	L *zerolog.Logger
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
		L: &logger,
	}
}

func (l *ZeroLogLogger) formatMessage(message interface{}) string {
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
func (l *ZeroLogLogger) Debug(message interface{}, args ...interface{}) {
	mf := l.formatMessage(message)
	l.log(l.L.Debug(), mf, args...)
}

// Info -.
func (l *ZeroLogLogger) Info(message interface{}, args ...interface{}) {
	mf := l.formatMessage(message)
	l.log(l.L.Info(), mf, args...)
}

// Warn -.
func (l *ZeroLogLogger) Warn(message interface{}, args ...interface{}) {
	mf := l.formatMessage(message)
	l.log(l.L.Warn(), mf, args...)
}

// Error -.
func (l *ZeroLogLogger) Error(message interface{}, args ...interface{}) {
	mf := l.formatMessage(message)
	l.log(l.L.Error(), mf, args...)
}

// Fatal -.
func (l *ZeroLogLogger) Fatal(message interface{}, args ...interface{}) {
	mf := l.formatMessage(message)
	l.log(l.L.Fatal(), mf, args...)

	os.Exit(1)
}

func (l *ZeroLogLogger) log(e *zerolog.Event, m string, args ...interface{}) {
	if len(args) == 0 {
		e.Msg(m)
	} else {
		e.Msgf(m, args...)
	}
}
