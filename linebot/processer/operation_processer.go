package processer

import (
	"fmt"
	"linebot/applicationerror"
	"linebot/dao"
	"linebot/entity"
	"linebot/logger"
	"linebot/transfer"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var isOperating bool = false

func HandleEvents(bot *linebot.Client, events []*linebot.Event) error {
	g := dao.InitDB()
	sqlDB, _ := g.DB.DB()
	defer sqlDB.Close()
	validEvents, notActiveUserEvents := validateEvent(events, g)
	// (b-3)返却処理
	returnReplyToNotValidUsers(notActiveUserEvents, bot)

	if len(validEvents) == 0 {
		return nil
	}

	// 後続処理
	userOpMap, masterOperation := margeEvents(validEvents, g)

	result, err := handleMasterOperation(masterOperation)

	if err != nil {
		handleErrorResponse(userOpMap, bot, err, g)
		return err
	}

	handleResponse(userOpMap, result, bot, g)

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
	logger.Info(&logger.LBIF050001, replyToken, resText)
	_, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(resText)).Do()
	if err != nil {
		logger.WarnWithStackTrace(err, &logger.LBWR050001, replyToken, resText)
	}
	return err
}

// ひとつのWebHookに含まれるEventをマージする
// ユーザ単位のマージ結果と、全体のマージ結果を返却
func margeEvents(events []*linebot.Event, g *dao.GormDB) (map[string]entity.Operation, entity.OperationType) {
	// key: lineId
	userOperations := map[string]entity.Operation{}
	// defaultではCheckを詰めておく
	lastOperation := entity.Check
	for _, e := range events {
		// validEventsでeventがTextMessageであること
		// MessageTextがopen or close or checkであることは担保済み
		op := entity.ConvertEventToOperatin(e)

		// 鍵の状態を変更する操作を優先するためCheck以外の場合は全体の操作を上書き
		if op.Operation != entity.Check {
			lastOperation = op.Operation
		}

		before, exist := userOperations[op.UserId]
		if exist {
			if op.Operation == entity.Check {
				// 非初操作かつ、Checkの場合、CheckをMergedとして記録
				g.InsertOperationHistory(op.UserId, op.Operation, entity.Merged)
			} else {
				// 非初操作かつ、Check以外の場合、前回の操作をMergedとして記録、かつ、操作を上書き
				g.UpdateOperationHistoryByOperationId(before.OperationId, entity.Merged)
				record := g.InsertOperationHistory(op.UserId, op.Operation, entity.Operating)
				op.OperationId = *record.OperationId
				userOperations[op.UserId] = *op
			}
		} else {
			record := g.InsertOperationHistory(op.UserId, op.Operation, entity.Operating)
			op.OperationId = *record.OperationId
			userOperations[op.UserId] = *op
		}
	}
	for _, o := range userOperations {
		logger.Debug(fmt.Sprintf("UserId : %s, OperationId : %d", o.UserId, o.OperationId))
	}
	return userOperations, lastOperation
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

func returnReplyToNotValidUsers(target []*linebot.Event, bot *linebot.Client) {
	for _, e := range target {
		logger.Info(&logger.LBIF020001, e.Source.UserID)
		reply(fmt.Sprintf("無効なユーザだよ。↓の文字列を管理者に送って。\n「%s」", e.Source.UserID), e.ReplyToken, bot)
	}
}

func handleResponse(ops map[string]entity.Operation, result entity.KeyServerResponse, bot *linebot.Client, g *dao.GormDB) {
	for _, o := range ops {
		if result.OperationStatus == "another" {
			go g.UpdateOperationHistoryWithErrorByOperationId(o.OperationId, entity.InOperatingError)
			reply("＝＝＝他操作の処理中＝＝＝", o.ReplyToken, bot)
			continue
		}
		// Check要求なら結果をそのまま返す
		go g.UpdateOperationHistoryByOperationId(o.OperationId, entity.Success)
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
}

func handleErrorResponse(ops map[string]entity.Operation, bot *linebot.Client, err error, g *dao.GormDB) {
	errorResponse := "エラーが起きてる！\nこのメッセージ見たらなるちゃんに「鍵のエラーハンドリングバグってるよ!」と連絡！"
	for _, o := range ops {
		switch err {
		case &applicationerror.ConnectionError:
			go g.UpdateOperationHistoryWithErrorByOperationId(o.OperationId, entity.KeyServerConnectionError)
			errorResponse = "＜＜鍵サーバとの通信に失敗した＞＞\nなるちゃんに連絡して!"
		case &applicationerror.ResponseParseError:
			go g.UpdateOperationHistoryWithErrorByOperationId(o.OperationId, entity.KeyServerConnectionError)
			errorResponse = fmt.Sprintf("！！！何が起きたか分からない！！！\nなるちゃんに↓これと一緒に至急連絡\n%s", applicationerror.ResponseParseError.Code)
		}
		reply(errorResponse, o.ReplyToken, bot)
	}
}
