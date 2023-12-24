package entity

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type OperationType int

const (
	Open  OperationType = 0
	Close OperationType = 1
	Check OperationType = 2
)

type Operation struct {
	OperationId int
	UserId      string
	Operation   OperationType
	ReplyToken  string
}

func TextToOperation(text string) OperationType {
	switch text {
	case "open":
		return Open
	case "close":
		return Close
	case "check":
		return Check
	default:
		return -1
	}
}

func ConvertEventToOperatin(event *linebot.Event) *Operation {
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		return &Operation{OperationId: -1, UserId: event.Source.UserID, Operation: TextToOperation(message.Text), ReplyToken: event.ReplyToken}
	default:
		return nil
	}
}
