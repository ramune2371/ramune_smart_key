package processer

import (
	"linebot/dao"
	"linebot/logger"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// WebHook Eventsの中身を検証
// 正常なeventのみ詰め直した配列と、無効なユーザからのEventの配列を返却
func validateEvent(events []*linebot.Event) ([]*linebot.Event, []*linebot.Event) {

	var verifiedEvent []*linebot.Event
	var notActiveUserEvent []*linebot.Event
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

		lineId := e.Source.UserID
		// (b-3)検証
		if !verifyUser(lineId) {
			logger.Debug(("not valid user"))
			notActiveUserEvent = append(notActiveUserEvent, e)
			dao.UpsertInvalidUser(lineId)
			continue
		}
		verifiedEvent = append(verifiedEvent, e)
		dao.UpdateUserLastAccess(lineId)
	}
	return verifiedEvent, notActiveUserEvent
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
