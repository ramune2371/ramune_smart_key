package transfer

import (
	"linebot/logger"
	"linebot/props"
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var bot *linebot.Client

func InitLineBot() {
	lineBot, err := linebot.New(props.ChannelSecret, props.ChannelToken)

	if err != nil {
		logger.FatalWithStackTrace(err, &logger.LBFT040004)
		panic(err)
	}
	bot = lineBot
}

func ParseLineRequest(r *http.Request) ([]*linebot.Event, error) {
	return bot.ParseRequest(r)
}

func ReplyToToken(resText, replyToken string) error {
	logger.Info(&logger.LBIF050001, replyToken, resText)
	_, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(resText)).Do()
	if err != nil {
		logger.WarnWithStackTrace(err, &logger.LBWR050001, replyToken, resText)
	}
	return err
}
