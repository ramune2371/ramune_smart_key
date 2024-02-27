package line

type LineTransfer interface {
	ReplyToToken(resText, replyToken string) error
}
