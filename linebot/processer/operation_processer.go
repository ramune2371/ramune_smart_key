package processer

import (
	"fmt"
	"linebot/applicationerror"
	"linebot/entity"
	"linebot/logger"
	"linebot/transfer"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var isOperating bool = false

func HandleEvents(bot *linebot.Client, events []*linebot.Event) error {
	validEvents, notActiveUserEvents := validateEvent(events)
	// (b-3)返却処理
	for _, e := range notActiveUserEvents {
		reply(fmt.Sprintf("無効なユーザだよ。↓の文字列を管理者に送って。\n「%s」", e.Source.UserID), e.ReplyToken, bot)
	}

	if len(validEvents) == 0 {
		logger.Debug("empty valid events")
		return nil
	}

	// 後続処理
	userOpMap, masterOperation := margeEvents(validEvents)
	result, err := handleMasterOperation(masterOperation)

	if err != nil {
		errorResponse := "エラーが起きてる！\nこのメッセージ見たらなるちゃんに「鍵のエラーハンドリングバグってるよ!」と連絡！"
		switch err {
		case &applicationerror.ConnectionError:
			errorResponse = "＜＜鍵サーバとの通信に失敗した＞＞\nなるちゃんに連絡して!"
		case &applicationerror.ResponseParseError:
			errorResponse = fmt.Sprintf("！！！何が起きたか分からない！！！\nなるちゃんに↓これと一緒に至急連絡\n%s", applicationerror.ResponseParseError.Code)
		}
		for _, o := range userOpMap {
			reply(errorResponse, o.ReplyToken, bot)
		}
		return err
	}

	for _, o := range userOpMap {
		if result.OperationStatus == "another" {
			reply("ーー操作中ーー", o.ReplyToken, bot)
			continue
		}
		// Check要求なら結果をそのまま返す
		if o.Operation == entity.Check {
			replyCheckResult(o.ReplyToken, result.KeyStatus, bot)
		} else {
			if o.Operation == entity.Open && result.KeyStatus == "True" {
				reply("→鍵開けたで", o.ReplyToken, bot)
			} else if o.Operation == entity.Close && result.KeyStatus == "False" {
				reply("→鍵閉めたで", o.ReplyToken, bot)
			} else if result.KeyStatus == "True" {
				reply("→誰かが開けたよ", o.ReplyToken, bot)
			} else {
				reply("→誰かが閉めたよ", o.ReplyToken, bot)
			}
		}
	}

	return nil
}

func replyCheckResult(replyToken string, result string, bot *linebot.Client) {
	if result == "True" {
		reply("あいてるよ", replyToken, bot)
	} else {
		reply("しまってるよ", replyToken, bot)
	}
}

func reply(resText, replyToken string, bot *linebot.Client) error {
	_, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(resText)).Do()
	return err
}

// ひとつのWebHookに含まれるEventをマージする
// ユーザ単位のマージ結果と、全体のマージ結果を返却
func margeEvents(events []*linebot.Event) (map[string]entity.Operation, entity.OperationType) {
	ret := map[string]entity.Operation{}
	// defaultではCheckを詰めておく
	lastOperation := entity.Check
	for _, e := range events {
		// validEventsでeventがTextMessageであること
		// MessageTextがopen or close or checkであること
		// は担保済み
		op := entity.ConvertEventToOperatin(e)

		// 鍵の状態を変更する操作を優先するためCheck以外の場合は全体の操作を上書き
		if op.Operation != entity.Check {
			lastOperation = op.Operation
		}

		// userの初操作または、操作がCheck以外ならユーザ操作を上書き
		_, ok := ret[op.UserId]
		if !ok || op.Operation != entity.Check {
			ret[op.UserId] = *op
		}
	}
	return ret, lastOperation
}

// 鍵操作の実行
func handleMasterOperation(operation entity.OperationType) (entity.KeyServerResponse, error) {

	if isOperating {
		ret := entity.KeyServerResponse{KeyStatus: "unknown", OperationStatus: "another"}
		return ret, nil
	}

	var ret entity.KeyServerResponse
	var err error
	isOperating = true
	switch operation {
	case entity.Open:
		ret, err = transfer.Open()
	case entity.Close:
		ret, err = transfer.Close()
	case entity.Check:
		ret, err = transfer.Check()
	default:
	}
	isOperating = false
	return ret, err
}
