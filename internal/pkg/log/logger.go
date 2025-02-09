package log

import (
	"os"

	"github.com/rs/zerolog"
)

const projectName = "go-project-layout"

type Logger interface {
	Trace() *zerolog.Event
	Debug() *zerolog.Event
	Info() *zerolog.Event
	Warn() *zerolog.Event
	Error() *zerolog.Event
}

func New(level string, serviceName string) *zerolog.Logger {
	zerolog.TimeFieldFormat = "02-01-2006 15:04:05"
	zerolog.LevelWarnValue = "warning"
	zerolog.MessageFieldName = "description"

	var isNotParsed bool
	var zerologLevel zerolog.Level
	switch level {
	case "trace":
		zerologLevel = zerolog.TraceLevel
	case "debug":
		zerologLevel = zerolog.DebugLevel
	case "info":
		zerologLevel = zerolog.InfoLevel
	case "warning":
		zerologLevel = zerolog.WarnLevel
	case "error":
		zerologLevel = zerolog.ErrorLevel
	default:
		isNotParsed = true
		zerologLevel = zerolog.InfoLevel
	}

	logger := zerolog.New(os.Stderr).
		Level(zerologLevel).
		With().
		Timestamp().
		Str("project", projectName).
		Str("service", serviceName).
		Logger()

	if isNotParsed {
		logger.Warn().Msgf("failed to parse loglevel: %s; level `info` set as default", level)
	}

	return &logger
}
