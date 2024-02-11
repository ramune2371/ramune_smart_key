package line

import (
	"linebot/logger"
	"linebot/props"
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type LineTransferImpl struct {
	bot *linebot.Client
}

func (ltImpl *LineTransferImpl) InitLineBot() {
	lineBot, err := linebot.New(props.ChannelSecret, props.ChannelToken)

	if err != nil {
		logger.FatalWithStackTrace(err, &logger.LBFT040004)
		panic(err)
	}
	ltImpl.bot = lineBot
}

func (ltImpl LineTransferImpl) ParseLineRequest(r *http.Request) ([]*linebot.Event, error) {
	return ltImpl.bot.ParseRequest(r)
}

func (ltImpl LineTransferImpl) ReplyToToken(resText, replyToken string) error {
	logger.Info(&logger.LBIF050001, replyToken, resText)
	_, err := ltImpl.bot.ReplyMessage(replyToken, linebot.NewTextMessage(resText)).Do()
	if err != nil {
		logger.WarnWithStackTrace(err, &logger.LBWR050001, replyToken, resText)
	}
	return err
}
