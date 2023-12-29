package processor

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

func HandleEvents(events []*linebot.Event) {
	validEvents, notActiveUserEvents := validateEvent(events)
	// (b-3)返却処理
	replyToNotValidUsers(notActiveUserEvents)

	if len(validEvents) == 0 {
		return
	}

	// 後続処理
	userOpMap, masterOperation := margeEvents(validEvents)

	if isOperating {
		for _, op := range userOpMap {
			replyInOperatingError(op)
		}
		return
	}

	result, err := handleMasterOperation(masterOperation)
	if err != nil {
		handleErrorResponse(userOpMap, err)
	}
	handleResponse(userOpMap, result)
}

func replyCheckResult(replyToken string, result string) {
	if result == "True" {
		transfer.ReplyToToken("あいてるよ", replyToken)
	} else {
		transfer.ReplyToToken("しまってるよ", replyToken)
	}
}

// ひとつのWebHookに含まれるEventをマージする
// ユーザ単位のマージ結果と、全体のマージ結果を返却
func margeEvents(events []*linebot.Event) (map[string]entity.Operation, entity.OperationType) {
	// key: lineId
	userOperations := map[string]entity.Operation{}
	// defaultではCheckを詰めておく
	lastOperation := entity.Check
	for _, e := range events {
		// validEventsでeventがTextMessageであること
		// MessageTextがopen or close or checkであることは担保済み
		op := entity.ConvertEventToOperation(e)

		// 鍵の状態を変更する操作を優先するためCheck以外の場合は全体の操作を上書き
		if op.Operation != entity.Check {
			lastOperation = op.Operation
		}

		before, exist := userOperations[op.UserId]
		if exist {
			if op.Operation == entity.Check {
				// 非初操作かつ、Checkの場合、CheckをMergedとして記録
				dao.InsertOperationHistory(op.UserId, op.Operation, entity.Merged)
			} else {
				// 非初操作かつ、Check以外の場合、前回の操作をMergedとして記録、かつ、操作を上書き
				dao.UpdateOperationHistoryByOperationId(before.OperationId, entity.Merged)
				record := dao.InsertOperationHistory(op.UserId, op.Operation, entity.Operating)
				op.OperationId = *record.OperationId
				userOperations[op.UserId] = *op
			}
		} else {
			record := dao.InsertOperationHistory(op.UserId, op.Operation, entity.Operating)
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

	isOperating = true
	var ret entity.KeyServerResponse
	var err error
	switch operation {
	case entity.Open:
		ret, err = transfer.OpenKey()
	case entity.Close:
		ret, err = transfer.CloseKey()
	case entity.Check:
		ret, err = transfer.CheckKey()
	default:
	}
	isOperating = false
	return ret, err
}

func replyToNotValidUsers(target []*linebot.Event) {
	for _, e := range target {
		logger.Info(&logger.LBIF020001, e.Source.UserID)
		transfer.ReplyToToken(fmt.Sprintf("無効なユーザだよ。↓の文字列を管理者に送って。\n「%s」", e.Source.UserID), e.ReplyToken)
	}
}

func replyInOperatingError(op entity.Operation) {
	go dao.UpdateOperationHistoryWithErrorByOperationId(op.OperationId, entity.InOperatingError)
	transfer.ReplyToToken("＝＝＝他操作の処理中＝＝＝", op.ReplyToken)
}

func handleResponse(ops map[string]entity.Operation, result entity.KeyServerResponse) {
	for _, o := range ops {
		if result.OperationStatus == "another" {
			replyInOperatingError(o)
			continue
		}
		go dao.UpdateOperationHistoryByOperationId(o.OperationId, entity.Success)
		// Check要求なら結果をそのまま返す
		if o.Operation == entity.Check {
			replyCheckResult(o.ReplyToken, result.KeyStatus)
		} else {
			if o.Operation == entity.Open && result.KeyStatus == "True" {
				transfer.ReplyToToken("→鍵開けたで", o.ReplyToken)
			} else if o.Operation == entity.Close && result.KeyStatus == "False" {
				transfer.ReplyToToken("→鍵閉めたで", o.ReplyToken)
			} else if result.KeyStatus == "True" {
				transfer.ReplyToToken("→誰かが開けたよ", o.ReplyToken)
			} else {
				transfer.ReplyToToken("→誰かが閉めたよ", o.ReplyToken)
			}
		}
	}
}

func handleErrorResponse(ops map[string]entity.Operation, err error) {
	errorResponse := "エラーが起きてる！\nこのメッセージ見たらなるちゃんに「鍵のエラーハンドリングバグってるよ!」と連絡！"
	for _, o := range ops {
		switch err {
		case &applicationerror.ConnectionError:
			go dao.UpdateOperationHistoryWithErrorByOperationId(o.OperationId, entity.KeyServerConnectionError)
			errorResponse = "＜＜鍵サーバとの通信に失敗した＞＞\nなるちゃんに連絡して!"
		case &applicationerror.ResponseParseError:
			go dao.UpdateOperationHistoryWithErrorByOperationId(o.OperationId, entity.KeyServerConnectionError)
			errorResponse = fmt.Sprintf("！！！何が起きたか分からない！！！\nなるちゃんに↓これと一緒に至急連絡\n%s", applicationerror.ResponseParseError.Code)
		}
		transfer.ReplyToToken(errorResponse, o.ReplyToken)
	}
}
