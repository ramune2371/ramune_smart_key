package entity

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type OperationType string

const (
	Open  OperationType = "open"
	Close OperationType = "close"
	Check OperationType = "check"
)

type Operation struct {
	UserId     string
	Operation  OperationType
	ReplyToken string
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
		return ""
	}
}

func ConvertEventToOperatin(event *linebot.Event) *Operation {
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		return &Operation{UserId: event.Source.UserID, Operation: TextToOperation(message.Text), ReplyToken: event.ReplyToken}
	default:
		return nil
	}
}
