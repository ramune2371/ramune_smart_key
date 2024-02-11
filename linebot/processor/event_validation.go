package processor

import (
	"linebot/dao/user_info"
	"linebot/entity"
	"linebot/logger"
	"linebot/security"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type EventValidator struct {
	UserInfoDao user_info.UserInfoDao
	Encryptor   security.Encryptor
}

// WebHook Eventsの中身を検証
// 正常なentity.Operationの配列および、無効なユーザのentity.Operationの配列を返却
func (eventValidation EventValidator) validateEvent(events []*linebot.Event) ([]*entity.Operation, []*entity.Operation) {

	allOperation := eventValidation.checkMessage(events)
	var verifiedOperation []*entity.Operation
	var notActiveUserOperation []*entity.Operation

	for _, op := range allOperation {
		lineId := op.UserId
		// (b-3)検証
		if !eventValidation.verifyUser(lineId) {
			logger.Debug(("not valid user"))
			notActiveUserOperation = append(notActiveUserOperation, op)
			continue
		}
		verifiedOperation = append(verifiedOperation, op)
	}
	return verifiedOperation, notActiveUserOperation
}

// LineBot EventがTextMessageかを検証
func (eventValidation EventValidator) isTextMessage(e *linebot.Event) bool {

	tm, ok := e.Message.(*linebot.TextMessage)
	if !ok {
		return false
	} else {
		return tm.Text != ""
	}
}

// TextMessageの中身が許可されたものかを検証
func (eventValidation EventValidator) verifyMessageText(text string) bool {

	return text != "" && (text == "open" || text == "close" || text == "check")
}

// Userが有効かを検証
func (eventValidation EventValidator) verifyUser(userId string) bool {

	user := eventValidation.UserInfoDao.GetUserByLineId(userId)
	if user == nil {
		return false
	}
	return user.Active
}

// Eventの配列のうち、Operationとして扱えるもののみを変換した配列を返却
func (ev EventValidator) checkMessage(events []*linebot.Event) []*entity.Operation {
	var allOperation []*entity.Operation
	for _, e := range events {
		//(b-1)検証
		if !ev.isTextMessage(e) {
			logger.Debug("not Text Message")
			continue
		}

		// 前段で型検証は済んでいるので型変換チェックはしない
		// (b-2)検証
		if !ev.verifyMessageText(e.Message.(*linebot.TextMessage).Text) {
			logger.Debug("not valid message")
			continue
		}

		allOperation = append(allOperation, EventConverterImpl{ev.Encryptor}.ConvertEventToEncryptedOperation(e))
	}

	return allOperation
}
