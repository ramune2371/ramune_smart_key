package processor

import (
	"fmt"
	"linebot/applicationerror"
	"linebot/dao/operation_history"
	"linebot/dao/user_info"
	"linebot/entity"
	"linebot/entity/message"
	"linebot/logger"
	"linebot/security"
	"linebot/transfer/key_server"
	"linebot/transfer/line"
	"sync"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type OperationProcessor struct {
	ohDao      operation_history.OperationHistoryDao
	uiDao      user_info.UserInfoDao
	lTransfer  line.LineTransfer
	ksTransfer key_server.KeyServerTransfer
	encryptor  security.Encryptor
	// 操作中判定フラグ(true:操作中, false:非操作中)
	isOperating bool
	// 操作中判定フラグのスレッドセーフ処理用
	mutex sync.Mutex
}

// OperationProcessorコンストラクタ
func NewOperationProcessor(ohDao operation_history.OperationHistoryDao, uiDao user_info.UserInfoDao, lTransfer line.LineTransfer, ksTransfer key_server.KeyServerTransfer, encryptor security.Encryptor) *OperationProcessor {
	op := &OperationProcessor{}
	op.isOperating = false
	op.ohDao = ohDao
	op.uiDao = uiDao
	op.lTransfer = lTransfer
	op.ksTransfer = ksTransfer
	op.encryptor = encryptor
	return op
}

/*
isOperationを設定

args: 設定する値(true:操作中, false:非操作中)
*/
func (op *OperationProcessor) SetIsOperating(value bool) {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	op.isOperating = value
}

/*
isOperationの値を取得

return: true:操作中, false:非操作中
*/
func (op *OperationProcessor) IsOperating() bool {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	return op.isOperating
}

/*
linebot.Eventを処理

args: events 処理対象のlinebot.Eventのポインタ配列
*/
func (op *OperationProcessor) HandleEvents(events []*linebot.Event) {
	validator := NewEventValidatorImpl(op.uiDao, op.encryptor)
	validEvents, invalidEvents := validator.ValidateEvent(events)

	if len(invalidEvents) != 0 {
		op.handleInvalidUserEvents(invalidEvents)
	}

	if len(validEvents) != 0 {
		op.handleValidUserEvents(validEvents)
	}

}

/*
無効なユーザからの操作を処理

ユーザの記録及び、ラインへの返答を実施

args: invalidUserEvents 処理対象のOperationのポインタ配列
*/
func (op *OperationProcessor) handleInvalidUserEvents(invalidUserEvents []*entity.Operation) {
	// 不正ユーザの記録
	for _, invalidEvent := range invalidUserEvents {
		op.uiDao.UpsertInvalidUser(invalidEvent.UserId)
	}

	// (b-3)返却処理
	op.replyToNotValidUsers(invalidUserEvents)
}

/*
有効なユーザからの操作を処理

ユーザの記録及び、ラインへの返答を実施

args: validUserEvents 処理対象のOperationのポインタ配列
*/
func (op *OperationProcessor) handleValidUserEvents(validEvents []*entity.Operation) {
	// 有効なユーザのアクセス記録
	for _, validEvent := range validEvents {
		op.uiDao.UpdateUserLastAccess(validEvent.UserId)
	}

	// 後続処理
	userOpMap, masterOperation := op.mergeEvents(validEvents)

	if op.IsOperating() {
		for _, operation := range userOpMap {
			op.replyInOperatingError(operation)
		}
		return
	}

	result, err := op.handleMasterOperation(masterOperation)
	if err != nil {
		op.handleKeyServerError(userOpMap, err)
		return
	}
	op.handleKeyServerResult(userOpMap, result)
}

/*
ひとつのWebHookに含まれるEventをマージする

args: operations 全ての操作要求
return: ユーザ単位の操作要求のマージ結果(key:UserId,value:Operation), 全体での操作要求のマージ結果
*/
func (opProcessor *OperationProcessor) mergeEvents(operations []*entity.Operation) (map[string]entity.Operation, entity.OperationType) {
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
				opProcessor.ohDao.InsertOperationHistory(op.UserId, op.Operation, entity.Merged)
			} else {
				// 非初操作かつ、Check以外の場合、前回の操作をMergedとして記録、かつ、操作を上書き
				opProcessor.ohDao.UpdateOperationHistoryByOperationId(before.OperationId, entity.Merged)
				record, _ := opProcessor.ohDao.InsertOperationHistory(op.UserId, op.Operation, entity.Operating)
				op.OperationId = *record.OperationId
				userOperations[op.UserId] = *op
			}
		} else {
			record, _ := opProcessor.ohDao.InsertOperationHistory(op.UserId, op.Operation, entity.Operating)
			op.OperationId = *record.OperationId
			userOperations[op.UserId] = *op
		}
	}
	for _, o := range userOperations {
		logger.Debug(fmt.Sprintf("UserId : %s, OperationId : %d", o.UserId, o.OperationId))
	}
	return userOperations, lastOperation
}

/*
鍵操作の実行

args: operation ユーザ全体の操作要求をマージした操作種別
return: 鍵サーバのレスポンス, error
*/
func (op *OperationProcessor) handleMasterOperation(operation entity.OperationType) (entity.KeyServerResponse, error) {

	op.SetIsOperating(true)
	defer op.SetIsOperating(false)

	var ret entity.KeyServerResponse
	var err error
	switch operation {
	case entity.Open:
		ret, err = op.ksTransfer.OpenKey()
	case entity.Close:
		ret, err = op.ksTransfer.CloseKey()
	case entity.Check:
		ret, err = op.ksTransfer.CheckKey()
	default:
		err = applicationerror.UnsupportedOperationError
	}
	return ret, err
}

/*
不正なユーザへの返信処理

args: target 返信対象のentity.Operationのポインタ配列
*/
func (opProcessor *OperationProcessor) replyToNotValidUsers(target []*entity.Operation) {
	for _, op := range target {
		logger.Info(logger.LBIF020001, op.UserId)
		msg := fmt.Sprintf("無効なユーザだよ。↓の文字列を管理者に送って。\n「%s」", op.UserId)
		opProcessor.lTransfer.ReplyToToken(msg, op.ReplyToken)
	}
}

/*
Check要求への応答処理

args: replyToken 返信対象のReplyToken, result 鍵の状態
*/
func (opProcessor *OperationProcessor) replyCheckResult(replyToken string, result entity.KeyStatus) {
	if result == entity.KeyStatusOpen {
		opProcessor.lTransfer.ReplyToToken(message.CHECK_OPEN, replyToken)
	} else {
		opProcessor.lTransfer.ReplyToToken(message.CHECK_CLOSE, replyToken)
	}
}

/*
操作中エラーのライン応答

args: op 応答対象のOperation
*/
func (opProcessor *OperationProcessor) replyInOperatingError(op entity.Operation) {
	go opProcessor.ohDao.UpdateOperationHistoryWithErrorByOperationId(op.OperationId, entity.InOperatingError)
	opProcessor.lTransfer.ReplyToToken(message.ANOTHER_OPERATING, op.ReplyToken)
}

/*
鍵サーバからの応答に応じてユーザからの操作要求へ返信

args: ops ユーザごとのマージ後の操作要求, result 鍵サーバからのレスポンス
*/
func (opProcessor *OperationProcessor) handleKeyServerResult(ops map[string]entity.Operation, result entity.KeyServerResponse) {
	for _, o := range ops {
		if result.OperationStatus == entity.OperationAnother {
			opProcessor.replyInOperatingError(o)
			continue
		}
		go opProcessor.ohDao.UpdateOperationHistoryByOperationId(o.OperationId, entity.Success)
		// Check要求なら結果をそのまま返す
		if o.Operation == entity.Check {
			opProcessor.replyCheckResult(o.ReplyToken, result.KeyStatus)
		} else {
			if o.Operation == entity.Open && result.KeyStatus == entity.KeyStatusOpen {
				opProcessor.lTransfer.ReplyToToken(message.SUCCESS_OPEN, o.ReplyToken)
			} else if o.Operation == entity.Close && result.KeyStatus == entity.KeyStatusClose {
				opProcessor.lTransfer.ReplyToToken(message.SUCCESS_CLOSE, o.ReplyToken)
			} else if result.KeyStatus == entity.KeyStatusOpen {
				opProcessor.lTransfer.ReplyToToken(message.ANOTHER_OPEN, o.ReplyToken)
			} else {
				opProcessor.lTransfer.ReplyToToken(message.ANOTHER_CLOSE, o.ReplyToken)
			}
		}
	}
}

/*
鍵サーバ接続時のエラーを処理

args: ops ユーザごとのマージ後の操作要求, err 鍵サーバ接続時のエラー
*/
func (opProcessor *OperationProcessor) handleKeyServerError(ops map[string]entity.Operation, err error) {
	errorResponse := "エラーが起きてる！\nこのメッセージ見たらなるちゃんに「鍵のエラーハンドリングバグってるよ!」と連絡！"
	for _, o := range ops {
		switch err {
		case applicationerror.ConnectionError:
			go opProcessor.ohDao.UpdateOperationHistoryWithErrorByOperationId(o.OperationId, entity.KeyServerConnectionError)
			errorResponse = message.CONNECTION_ERROR
		case applicationerror.ResponseParseError:
			go opProcessor.ohDao.UpdateOperationHistoryWithErrorByOperationId(o.OperationId, entity.KeyServerResponseError)
			errorResponse = fmt.Sprintf(message.APPLICATION_ERROR, applicationerror.ResponseParseError.Code)
		default:
			errorResponse = "不正な操作！どうやってここまで辿り着いた？？"
		}
		opProcessor.lTransfer.ReplyToToken(errorResponse, o.ReplyToken)
	}
}
