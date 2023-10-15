package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

func Request(req *http.Request){
  buf,err := io.ReadAll(req.Body)
  if err != nil {
    Fatal(fmt.Sprintf("Request Logging Failure: %s",err.Error()),"FT999999")
  }
  log.Info().Str("type","RequestLogger").RawJSON("header",toJson(req.Header)).Msg(string(buf))
  req.Body = io.NopCloser(bytes.NewBuffer(buf))
}

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

func toJson(target any) []byte{
  jb,err := json.Marshal(target)
  if err != nil {
    Fatal(fmt.Sprintf("Json Marshal Filure: %s",err.Error()),"FT909001")
  }
  return jb
}
