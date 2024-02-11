package line

import (
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type LineTransfer interface {
	ParseLineRequest(r *http.Request) ([]*linebot.Event, error)
	ReplyToToken(resText, replyToken string) error
}
