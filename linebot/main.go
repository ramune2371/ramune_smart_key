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


func main() {
  bot, err := linebot.New(props.ChannelSecret, props.ChannelToken)
  if err != nil {
    fmt.Println(err)
  }
  e := echo.New()
  e.HideBanner = true
  //e.Use(middleware.Logger())
  e.Use(middle.VerifyLineSignature)
  e.POST("/", func(c echo.Context) error {
    events, err := bot.ParseRequest(c.Request())
    if err != nil {
      return c.NoContent(http.StatusBadRequest)
    }
    for i,e := range events{
      switch message := e.Message.(type){
      case *linebot.TextMessage:
        logger.Debug(fmt.Sprintf("%d,%s",i,message.Text))
        processer.HandleRequest(message.Text,e.ReplyToken,bot)
        break
      default:
        continue
      }
      
    }
    return c.String(http.StatusOK, "")
  })
  e.Logger.Fatal(e.Start(":1323"))
}
