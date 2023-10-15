package middle

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"linebot/props"
  "ramune/modules/logger"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

/**
 RequestのHTTPヘッダーにあるLINEの署名を検証。
 検証NGの場合は400Bad Requestを返却
**/
func VerifyLineSignature(next echo.HandlerFunc) echo.HandlerFunc{
  return func(c echo.Context) error {
    logger.Info("Start Verify Line Signature","IF101001")

    // Reqeustの読み込み
    req := c.Request()
    body, err := io.ReadAll(req.Body)
    if err != nil {
      logger.Warn(err.Error(),"WR101001")
      return c.NoContent(http.StatusInternalServerError)
    }

    // 署名の検証処理
    decoded, err := base64.StdEncoding.DecodeString(req.Header.Get("x-line-signature"))
    if err != nil {
      logger.Warn(err.Error(),"WR101002")
      return c.NoContent(http.StatusInternalServerError)
    }
    hash := hmac.New(sha256.New, []byte(props.ChannelSecret))
    hash.Write(body)
    if !hmac.Equal(decoded,hash.Sum(nil)) {
      logger.Warn("Signature is not verify","WR101003")
      return c.NoContent(http.StatusBadRequest)
    }

    // 後続でRequestBodyを処理するために詰め直し
    c.Request().Body = io.NopCloser(bytes.NewBuffer(body))
    logger.Info("End Verify Line Signature","IF101002")
    return next(c)
  }
}

func VerifyLineWebHookEvent(e linebot.Event) bool {
  
  if e.Type != linebot.EventTypeMessage || e.Message.Type()!=linebot.MessageTypeText {
    return false
  }
  return true
}
