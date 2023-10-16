package processer

import (
	"fmt"
	"linebot/dao"
	"linebot/logger"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// WebHook Eventsの中身を検証
// 正常なeventのみ詰め直した配列と、無効なユーザからのEventの配列を返却
func validateEvent(events []*linebot.Event) ([]*linebot.Event,[]*linebot.Event){

  var verifiedEvent []*linebot.Event
  var notActiveUserEvent []*linebot.Event
  for _,e := range events {
    //(b-1)検証
    if !isTextMessage(e) {
      continue
    }

    // 前段で型検証は済んでいるので型変換チェックはしない
    // (b-2)検証
    if !verifyMessageText(e.Message.(*linebot.TextMessage).Text) {
      continue
    }

    // (b-3)検証
    if !verifyUser(e.Source.UserID) {
      notActiveUserEvent = append(notActiveUserEvent, e)
      continue
    }
    verifiedEvent = append(verifiedEvent,e)
  }
  return verifiedEvent,notActiveUserEvent
}

// LineBot EventがTextMessageかを検証
func isTextMessage(e *linebot.Event)bool{

  _,ok := e.Message.(*linebot.TextMessage)
  return ok
}

// TextMessageの中身が許可されたものかを検証
func verifyMessageText(text string) bool {

  return text=="open"||text=="close"||text=="check"
}

// Userが有効かを検証
func verifyUser(userId string) bool {

  name,active := dao.GetUserActive(userId)
  logger.Debug(fmt.Sprintf("Request User id = %s, name = %s",userId,name))
  return active
}

// LINEからのリクエスト内容に応じて処理を実行
func HandleEvent(text,replyToken string,bot *linebot.Client) error {
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

func HandleEvents(bot *linebot.Client, events []*linebot.Event) error {
  validEvents,notActiveUserEvents := validateEvent(events)
  // (b-3)返却処理
  for _,e := range notActiveUserEvents {
    reply(fmt.Sprintf("無効なユーザだよ。↓の文字列を管理者に送って。\n「%s」",e.Source.UserID),e.ReplyToken,bot)
  }

  // 後続処理
  for i,e := range validEvents {
    switch message := e.Message.(type){
    case *linebot.TextMessage:
      logger.Debug(fmt.Sprintf("%d,%s",i,message.Text))
      HandleEvent(message.Text,e.ReplyToken,bot)
      break
    default:
      // Text Message以外は無視
      continue
    }
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
