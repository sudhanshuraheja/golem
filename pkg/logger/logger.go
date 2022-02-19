package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// Logger ...
type Logger struct {
	logger *zerolog.Logger
}

// NewLogger ...
func NewLogger(logLevel int, prettyPrint bool) *Logger {
	var zeroLogLevel zerolog.Level
	switch logLevel {
	case -1:
		zeroLogLevel = zerolog.TraceLevel
	case 0:
		zeroLogLevel = zerolog.DebugLevel
	case 1:
		zeroLogLevel = zerolog.InfoLevel
	case 2:
		zeroLogLevel = zerolog.WarnLevel
	case 3:
		zeroLogLevel = zerolog.ErrorLevel
	case 4:
		zeroLogLevel = zerolog.FatalLevel
	case 5:
		zeroLogLevel = zerolog.PanicLevel
	}

	zerolog.SetGlobalLevel(zeroLogLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	logger := zerolog.New(os.Stderr).
		With().
		Timestamp().
		Logger()

	if prettyPrint {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	return &Logger{logger: &logger}
}

// Get ...
func (l *Logger) Get() zerolog.Logger {
	return *l.logger
}

// Debug starts a new message with debug level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Debug(where string) *zerolog.Event {
	return l.logger.Debug().Str("z", where)
}

// Info starts a new message with info level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Info(where string) *zerolog.Event {
	return l.logger.Info().Str("z", where)
}

// Warn starts a new message with warn level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Warn(where string) *zerolog.Event {
	return l.logger.Warn().Str("z", where)
}

// Error starts a new message with error level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Error(where string) *zerolog.Event {
	return l.logger.Error().Str("z", where)
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the Msg method.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Fatal(where string) *zerolog.Event {
	return l.logger.Fatal().Str("z", where)
}

// Panic starts a new message with panic level. The message is also sent
// to the panic function.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Panic(where string) *zerolog.Event {
	return l.logger.Panic().Str("z", where)
}

// LogError ...
func (l *Logger) LogError(where string, err error) {
	l.Error(where).Msgf("%+v", err)
}
