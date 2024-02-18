package line

import (
	"linebot/applicationerror"
	"linebot/logger"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type LineTransferImpl struct {
	bot *linebot.Client
}

func NewLineTransfer(bot *linebot.Client) *LineTransferImpl {
	ltImpl := new(LineTransferImpl)
	ltImpl.bot = bot
	return ltImpl
}

func (ltImpl LineTransferImpl) ReplyToToken(resText, replyToken string) error {
	logger.Info(logger.LBIF050001, replyToken, resText)
	_, err := ltImpl.bot.ReplyMessage(replyToken, linebot.NewTextMessage(resText)).Do()
	if err != nil {
		logger.WarnWithStackTrace(err, applicationerror.ReplyError, logger.LBWR050001, replyToken, resText)
	}
	return err
}
