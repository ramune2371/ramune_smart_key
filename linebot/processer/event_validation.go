package processer

import (
	"linebot/dao"

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
		verifiedEvent = append(verifiedEvent, e)
	}
	return verifiedEvent, notActiveUserEvent
}

// LineBot EventがTextMessageかを検証
func isTextMessage(e *linebot.Event) bool {

	_, ok := e.Message.(*linebot.TextMessage)
	return ok
}

// TextMessageの中身が許可されたものかを検証
func verifyMessageText(text string) bool {

	return text == "open" || text == "close" || text == "check"
}

// Userが有効かを検証
func verifyUser(userId string) bool {

	user := dao.GetUserByLineId(userId)
	if user == nil {
		return false
	}
	return user.Active
}
