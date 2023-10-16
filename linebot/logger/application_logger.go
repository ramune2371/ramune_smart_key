package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

// TODO ログ設計

func Request(req *http.Request) {
  buf,err := io.ReadAll(req.Body)
  if err != nil {
    Fatal(fmt.Sprintf("Error at reading request Body err:%s",err.Error()),"FT999999")
  }

  log.Info().Str("type","RequestLogger").RawJSON("header",toJson(req.Header)).Msg(string(buf))

  // 後続でRequestの内容を読むために詰め直し
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
  ret,err := json.Marshal(target)
  if err != nil {
    Fatal(fmt.Sprintf("Error at Reading as Json err:%s",err.Error()),"FT999999")
    return nil
  }
  return ret
}


