package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

var log zerolog.Logger

func Init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logLevel := strings.ToLower(viper.GetString("LOG_LEVEL"))

	switch logLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

func GetLevel() zerolog.Level {
	return zerolog.GlobalLevel()
}

func Debug(msg string, fields ...interface{}) {
	event := log.Debug()
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			event = event.Interface(fields[i].(string), fields[i+1])
		}
	}
	event.Msg(msg)
}

func Info(msg string, fields ...interface{}) {
	event := log.Info()
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			event = event.Interface(fields[i].(string), fields[i+1])
		}
	}
	event.Msg(msg)
}

func Warn(msg string, fields ...interface{}) {
	event := log.Warn()
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			event = event.Interface(fields[i].(string), fields[i+1])
		}
	}
	event.Msg(msg)
}

func Error(msg string, fields ...interface{}) {
	event := log.Error()
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			event = event.Interface(fields[i].(string), fields[i+1])
		}
	}
	event.Msg(msg)
}
