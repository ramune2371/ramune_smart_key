package processor

import (
	"linebot/entity"
	"linebot/security"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type EventConverter interface {
	ConvertEventToEncryptedOperation(*linebot.Event) *entity.Operation
}

type EventConverterImpl struct {
	Encryptor security.Encryptor
}

func TextToOperation(text string) entity.OperationType {
	switch text {
	case "open":
		return entity.Open
	case "close":
		return entity.Close
	case "check":
		return entity.Check
	default:
		return entity.Unsupported
	}
}

// LINE Webhook Eventをentity.Operationに変換&LINE IDをソルト付きハッシュ化
func (ec EventConverterImpl) ConvertEventToEncryptedOperation(event *linebot.Event) *entity.Operation {
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		opType := TextToOperation(message.Text)
		if opType == entity.Unsupported {
			return nil
		}
		return &entity.Operation{
			OperationId: -1,
			UserId:      ec.Encryptor.SaltHash(event.Source.UserID),
			Operation:   opType,
			ReplyToken:  event.ReplyToken,
		}
	default:
		return nil
	}
}
