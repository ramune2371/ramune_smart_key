package logger

import (
	"github.com/rs/zerolog/log"
)

func Debug(message string) {
  log.Debug().Str("type","DebugLogger").Msg(message)
}

func Info(message,id string) {
  log.Info().Str("type","ApplicationLogger").Str("id",id).Msg(message)
}

func Warn(message,id string) {
  log.Warn().Str("type","ApplicationLogger").Str("id",id).Msg(message)
}

func Error(message,id string) {
  log.Error().Str("type","ApplicationLogger").Str("id",id).Msg(message)
}

func Fatal(message,id string) {
  log.Fatal().Str("type","ApplicationLogger").Str("id",id).Msg(message)
}
