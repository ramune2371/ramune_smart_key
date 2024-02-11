package processor

import (
	"fmt"
	"linebot/applicationerror"
	"linebot/dao/operation_history"
	"linebot/dao/user_info"
	"linebot/entity"
	"linebot/logger"
	"linebot/security"
	"linebot/transfer/key_server"
	"linebot/transfer/line"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type OperationProcessor struct {
	OpHistoryDao      operation_history.OperationHistoryDao
	UserInfoDao       user_info.UserInfoDao
	LineTransfer      line.LineTransfer
	KeyServerTransfer key_server.KeyServerTransfer
	Encryptor         security.Encryptor
}

var isOperating bool = false

func (opProcessor OperationProcessor) HandleEvents(events []*linebot.Event) {
	validator := EventValidator{UserInfoDao: opProcessor.UserInfoDao, Encryptor: opProcessor.Encryptor}
	validEvents, notActiveUserEvents := validator.validateEvent(events)

	// 不正ユーザの記録
	for _, invalidEvent := range notActiveUserEvents {
		opProcessor.UserInfoDao.UpsertInvalidUser(invalidEvent.UserId)
	}

	// 有効なユーザのアクセス記録
	for _, validEvent := range validEvents {
		opProcessor.UserInfoDao.UpdateUserLastAccess(validEvent.UserId)
	}

	// (b-3)返却処理
	opProcessor.replyToNotValidUsers(notActiveUserEvents)

	if len(validEvents) == 0 {
		return
	}

	// 後続処理
	userOpMap, masterOperation := opProcessor.margeEvents(validEvents)

	if isOperating {
		for _, op := range userOpMap {
			opProcessor.replyInOperatingError(op)
		}
		return
	}

	result, err := opProcessor.handleMasterOperation(masterOperation)
	if err != nil {
		opProcessor.handleKeyServerError(userOpMap, err)
		return
	}
	opProcessor.handleKeyServerResult(userOpMap, result)
}

// ひとつのWebHookに含まれるEventをマージする
// ユーザ単位のマージ結果と、全体のマージ結果を返却
func (opProcessor OperationProcessor) margeEvents(operations []*entity.Operation) (map[string]entity.Operation, entity.OperationType) {
	// key: lineId
	userOperations := map[string]entity.Operation{}
	// defaultではCheckを詰めておく
	lastOperation := entity.Check
	for _, op := range operations {

		// 鍵の状態を変更する操作を優先するためCheck以外の場合は全体の操作を上書き
		if op.Operation != entity.Check {
			lastOperation = op.Operation
		}

		before, exist := userOperations[op.UserId]
		if exist {
			if op.Operation == entity.Check {
				// 非初操作かつ、Checkの場合、CheckをMergedとして記録
				opProcessor.OpHistoryDao.InsertOperationHistory(op.UserId, op.Operation, entity.Merged)
			} else {
				// 非初操作かつ、Check以外の場合、前回の操作をMergedとして記録、かつ、操作を上書き
				opProcessor.OpHistoryDao.UpdateOperationHistoryByOperationId(before.OperationId, entity.Merged)
				record := opProcessor.OpHistoryDao.InsertOperationHistory(op.UserId, op.Operation, entity.Operating)
				op.OperationId = *record.OperationId
				userOperations[op.UserId] = *op
			}
		} else {
			record := opProcessor.OpHistoryDao.InsertOperationHistory(op.UserId, op.Operation, entity.Operating)
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
func (opProcessor OperationProcessor) handleMasterOperation(operation entity.OperationType) (entity.KeyServerResponse, error) {

	isOperating = true
	var ret entity.KeyServerResponse
	var err error
	switch operation {
	case entity.Open:
		ret, err = opProcessor.KeyServerTransfer.OpenKey()
	case entity.Close:
		ret, err = opProcessor.KeyServerTransfer.CloseKey()
	case entity.Check:
		ret, err = opProcessor.KeyServerTransfer.CheckKey()
	default:
	}
	isOperating = false
	return ret, err
}

func (opProcessor OperationProcessor) replyToNotValidUsers(target []*entity.Operation) {
	for _, op := range target {
		logger.Info(&logger.LBIF020001, op.UserId)
		msg := fmt.Sprintf("無効なユーザだよ。↓の文字列を管理者に送って。\n「%s」", op.UserId)
		opProcessor.LineTransfer.ReplyToToken(msg, op.ReplyToken)
	}
}

func (opProcessor OperationProcessor) replyCheckResult(replyToken string, result entity.KeyStatus) {
	if result == entity.KeyStatusOpen {
		opProcessor.LineTransfer.ReplyToToken("あいてるよ", replyToken)
	} else {
		opProcessor.LineTransfer.ReplyToToken("しまってるよ", replyToken)
	}
}

func (opProcessor OperationProcessor) replyInOperatingError(op entity.Operation) {
	go opProcessor.OpHistoryDao.UpdateOperationHistoryWithErrorByOperationId(op.OperationId, entity.InOperatingError)
	opProcessor.LineTransfer.ReplyToToken("＝＝＝他操作の処理中＝＝＝", op.ReplyToken)
}

func (opProcessor OperationProcessor) handleKeyServerResult(ops map[string]entity.Operation, result entity.KeyServerResponse) {
	for _, o := range ops {
		if result.OperationStatus == entity.OperationAnother {
			opProcessor.replyInOperatingError(o)
			continue
		}
		go opProcessor.OpHistoryDao.UpdateOperationHistoryByOperationId(o.OperationId, entity.Success)
		// Check要求なら結果をそのまま返す
		if o.Operation == entity.Check {
			opProcessor.replyCheckResult(o.ReplyToken, result.KeyStatus)
		} else {
			if o.Operation == entity.Open && result.KeyStatus == entity.KeyStatusOpen {
				opProcessor.LineTransfer.ReplyToToken("→鍵開けたで", o.ReplyToken)
			} else if o.Operation == entity.Close && result.KeyStatus == entity.KeyStatusClose {
				opProcessor.LineTransfer.ReplyToToken("→鍵閉めたで", o.ReplyToken)
			} else if result.KeyStatus == entity.KeyStatusOpen {
				opProcessor.LineTransfer.ReplyToToken("→誰かが開けたよ", o.ReplyToken)
			} else {
				opProcessor.LineTransfer.ReplyToToken("→誰かが閉めたよ", o.ReplyToken)
			}
		}
	}
}

func (opProcessor OperationProcessor) handleKeyServerError(ops map[string]entity.Operation, err error) {
	errorResponse := "エラーが起きてる！\nこのメッセージ見たらなるちゃんに「鍵のエラーハンドリングバグってるよ!」と連絡！"
	for _, o := range ops {
		switch err {
		case &applicationerror.ConnectionError:
			go opProcessor.OpHistoryDao.UpdateOperationHistoryWithErrorByOperationId(o.OperationId, entity.KeyServerConnectionError)
			errorResponse = "＜＜鍵サーバとの通信に失敗した＞＞\nなるちゃんに連絡して!"
		case &applicationerror.ResponseParseError:
			go opProcessor.OpHistoryDao.UpdateOperationHistoryWithErrorByOperationId(o.OperationId, entity.KeyServerConnectionError)
			errorResponse = fmt.Sprintf("！！！何が起きたか分からない！！！\nなるちゃんに↓これと一緒に至急連絡\n%s", applicationerror.ResponseParseError.Code)
		}
		opProcessor.LineTransfer.ReplyToToken(errorResponse, o.ReplyToken)
	}
}
