package processer

import (
	"fmt"
	"linebot/dao"
	"linebot/entity"
	"linebot/transfer"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var isOperating bool = false

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

  user := dao.GetUserByLineId(userId)
  if user == nil {
    return false
  }
  return user.Active
}

func HandleEvents(bot *linebot.Client, events []*linebot.Event) error {
  validEvents,notActiveUserEvents := validateEvent(events)
  // (b-3)返却処理
  for _,e := range notActiveUserEvents {
    reply(fmt.Sprintf("無効なユーザだよ。↓の文字列を管理者に送って。\n「%s」",e.Source.UserID),e.ReplyToken,bot)
  }

  // 後続処理
  userOpMap,masterOperation := margeEvents(validEvents)
  result := handleMasterOperation(masterOperation)
  for _,o := range userOpMap {
    if result.OperationStatus == "another" {
      reply("ーー操作中ーー",o.ReplyToken,bot)
      continue
    }
    // Check要求なら結果をそのまま返す
    if o.Operation == entity.Check {
      replyCheckResult(o.ReplyToken,result.KeyStatus,bot)
    }else{
      if o.Operation == entity.Open && result.KeyStatus == "True" {
        reply("→鍵開けたで",o.ReplyToken,bot)
      }else if o.Operation == entity.Close && result.KeyStatus == "False" {
        reply("→鍵閉めたで",o.ReplyToken,bot)
      }else if result.KeyStatus == "True" {
        reply("！なんか知らんけど開いたわ！",o.ReplyToken,bot)
      }else {
        reply("！なんか知らんけど閉じたわ！",o.ReplyToken,bot)
      }
    } 
  }
  
  return nil
}

func replyCheckResult(replyToken string,result string,bot *linebot.Client){
  if result == "True" {
    reply("＜＜鍵開いてるで＞＞",replyToken,bot)
  }else{
    reply("＞＞鍵閉じてるで＜＜",replyToken,bot)
  }
}

func reply(resText,replyToken string,bot *linebot.Client) error {
  _,err := bot.ReplyMessage(replyToken,linebot.NewTextMessage(resText)).Do()
  return err
}

// ひとつのWebHookに含まれるEventをマージする
// ユーザ単位のマージ結果と、全体のマージ結果を返却
func margeEvents(events []*linebot.Event) (map[string]entity.Operation,entity.OperationType){
  ret := map[string]entity.Operation{}
  // defaultではCheckを詰めておく
  lastOperation := entity.Check
  for _,e := range events {
    // validEventsでeventがTextMessageであること
    // MessageTextがopen or close or checkであること
    // は担保済み
    op := entity.ConvertEventToOperatin(e)
    
    // 鍵の状態を変更する操作を優先するためCheck以外の場合は全体の操作を上書き
    if(op.Operation != entity.Check){
      lastOperation = op.Operation
    }

    // userの初操作または、操作がCheck以外ならユーザ操作を上書き
    _,ok := ret[op.UserId]
    if !ok || op.Operation != entity.Check{
      ret[op.UserId] = *op
    }
  }
  return ret,lastOperation
}

// true:open false:close
func handleMasterOperation(operation entity.OperationType) entity.KeyServerResponse {

  var ret entity.KeyServerResponse
  if isOperating {
    ret = entity.KeyServerResponse{KeyStatus:"unknown",OperationStatus:"another"}
    return ret
  }

  isOperating = true
  switch (operation){
  case entity.Open:
    ret = transfer.Open()
    break
  case entity.Close:
    ret = transfer.Close()
    break
  case entity.Check:
    ret = transfer.Check()
    break
  default:
  }
  isOperating = false
  return ret
}
