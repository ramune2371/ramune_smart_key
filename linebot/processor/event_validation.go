package processor

import (
	"linebot/dao"
	"linebot/entity"
	"linebot/logger"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// WebHook Eventsの中身を検証
// 正常なentity.Operationの配列および、無効なユーザのentity.Operationの配列を返却
func validateEvent(events []*linebot.Event) ([]*entity.Operation, []*entity.Operation) {

	allOperation := checkMessage(events)
	var verifiedOperation []*entity.Operation
	var notActiveUserOperation []*entity.Operation

	for _, op := range allOperation {
		lineId := op.UserId
		// (b-3)検証
		if !verifyUser(lineId) {
			logger.Debug(("not valid user"))
			notActiveUserOperation = append(notActiveUserOperation, op)
			dao.UpsertInvalidUser(lineId)
			continue
		}
		verifiedOperation = append(verifiedOperation, op)
		dao.UpdateUserLastAccess(lineId)
	}
	return verifiedOperation, notActiveUserOperation
}

// LineBot EventがTextMessageかを検証
func isTextMessage(e *linebot.Event) bool {

	tm, ok := e.Message.(*linebot.TextMessage)
	if !ok {
		return false
	} else {
		return tm.Text != ""
	}
}

// TextMessageの中身が許可されたものかを検証
func verifyMessageText(text string) bool {

	return text != "" && (text == "open" || text == "close" || text == "check")
}

// Userが有効かを検証
func verifyUser(userId string) bool {

	user := dao.GetUserByLineId(userId)
	if user == nil {
		return false
	}
	return user.Active
}

// Eventの配列のうち、Operationとして扱えるもののみを変換した配列を返却
func checkMessage(events []*linebot.Event) []*entity.Operation {
	var allOperation []*entity.Operation
	for _, e := range events {
		//(b-1)検証
		if !isTextMessage(e) {
			logger.Debug("not Text Message")
			continue
		}

		// 前段で型検証は済んでいるので型変換チェックはしない
		// (b-2)検証
		if !verifyMessageText(e.Message.(*linebot.TextMessage).Text) {
			logger.Debug("not valid message")
			continue
		}

		allOperation = append(allOperation, entity.ConvertEventToOperation(e))
	}

	return allOperation
}
