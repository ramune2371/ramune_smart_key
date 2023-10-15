package processer

import (
	"fmt"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// LINEからのリクエスト内容に応じて処理を実行
func HandleRequest(text,replyToken string,bot *linebot.Client) error {
  switch text {
  case "open":
    handleOpenOperation(bot)
    reply("開けるよー",replyToken,bot)
  case "close":
    handleCloseOperation(bot)
    reply("閉めたいの？",replyToken,bot)
    return nil
  case "check":
    handleCheckOperation(bot)
    reply("確認するね",replyToken,bot)
    return nil
  default:
    return fmt.Errorf("Un Supported Operation:%s",text)
  }
  return nil
}

func handleOpenOperation(bot *linebot.Client) error {
  return nil
}

func handleCloseOperation(bot *linebot.Client) error {
  return nil
}

func handleCheckOperation(bot *linebot.Client) (bool,error) {

  return true,nil
}

func reply(resText,replyToken string,bot *linebot.Client) error {
  _,err := bot.ReplyMessage(replyToken,linebot.NewTextMessage(resText)).Do()
  return err
}
