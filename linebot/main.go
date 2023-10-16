package main

import (
	"fmt"
	"linebot/middle"
	"linebot/processer"
	"linebot/props"
  "linebot/logger"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func koremo(){
  fmt.Println("jikko sareru?")
}

func requestLog(next echo.HandlerFunc) echo.HandlerFunc {
  return func(c echo.Context) error {
    logger.Request(c.Request())
    return next(c)
  }
}


func main() {
  bot, err := linebot.New(props.ChannelSecret, props.ChannelToken)
  if err != nil {
    fmt.Println(err)
  }
  e := echo.New()
  e.HideBanner = true
  e.Use(requestLog)
  e.Use(middle.VerifyLineSignature)
  e.POST("/", func(c echo.Context) error {
    events, err := bot.ParseRequest(c.Request())
    if err != nil {
      return c.NoContent(http.StatusBadRequest)
    }
    processer.HandleEvents(bot,events)

    return c.String(http.StatusOK, "")
  })
  e.Logger.Fatal(e.Start(":1323"))
}

