package processor

import (
	"linebot/entity"
	mock_security "linebot/security/mock"
	"linebot/testutil"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func TestTextToOperation(t *testing.T) {
	// テスト対象定義
	tests := []struct {
		description string
		text        string
		expect      entity.OperationType
	}{
		{
			description: "valid (open)",
			text:        "open",
			expect:      entity.Open,
		},
		{
			description: "valid (close)",
			text:        "close",
			expect:      entity.Close,
		},
		{
			description: "valid (check)",
			text:        "check",
			expect:      entity.Check,
		},
		{
			description: "invalid (not empty)",
			text:        "fizz",
			expect:      -1,
		},
		{
			description: "invalid (empty)",
			text:        "",
			expect:      -1,
		},
	}

	// テスト実行
	for _, test := range tests {
		ret := TextToOperation(test.text)
		if ret != test.expect {
			t.Errorf(testutil.INT_TEST_MSG_FMT, test.description, test.expect, ret)
		}
	}
}

func TestConvertEventToEncryptedOperation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEncryptor := mock_security.NewMockEncryptor(ctrl)
	mockEncryptor.EXPECT().SaltHash(gomock.Any()).DoAndReturn(func(value string) string { return value }).AnyTimes()

	ec := EventConverterImpl{mockEncryptor}

	userId := "userId"
	replyToken := "replyToken"

	// テスト対象定義
	tests := []struct {
		description string
		event       *linebot.Event
		expect      *entity.Operation
	}{
		{
			description: "Valid Text Event",
			event: &linebot.Event{
				Message: linebot.NewTextMessage("open"),
				Source: &linebot.EventSource{
					UserID: userId,
				},
				ReplyToken: replyToken,
			},
			expect: &entity.Operation{OperationId: -1, Operation: entity.Open, UserId: userId, ReplyToken: replyToken},
		},
		{
			description: "Invalid Text Event(empty)",
			event: &linebot.Event{
				Message: linebot.NewTextMessage(""),
				Source: &linebot.EventSource{
					UserID: userId,
				},
				ReplyToken: replyToken,
			},
			expect: nil,
		},
		{
			description: "Invalid Event(Audio Message)",
			event: &linebot.Event{
				Message: linebot.NewAudioMessage("test", 0),
				Source: &linebot.EventSource{
					UserID: userId,
				},
				ReplyToken: replyToken,
			},
			expect: nil,
		},
		{
			description: "Invalid Event(Flex Message)",
			event: &linebot.Event{
				Message: linebot.NewFlexMessage("test", nil),
				Source: &linebot.EventSource{
					UserID: userId,
				},
				ReplyToken: replyToken,
			},
			expect: nil,
		},
		{
			description: "Invalid Event(Image Message)",
			event: &linebot.Event{
				Message: linebot.NewImageMessage("test", "test"),
				Source: &linebot.EventSource{
					UserID: userId,
				},
				ReplyToken: replyToken,
			},
			expect: nil,
		},
		{
			description: "Invalid Event(Video Message)",
			event: &linebot.Event{
				Message: linebot.NewVideoMessage("", ""),
				Source: &linebot.EventSource{
					UserID: userId,
				},
				ReplyToken: replyToken,
			},
			expect: nil,
		},
		{
			description: "Invalid Event(Location Message)",
			event: &linebot.Event{
				Message: linebot.NewLocationMessage("title", "", 0, 0),
				Source: &linebot.EventSource{
					UserID: userId,
				},
				ReplyToken: replyToken,
			},
			expect: nil,
		},
		{
			description: "Invalid Event(Sticker Message)",
			event: &linebot.Event{
				Message: linebot.NewStickerMessage("", ""),
				Source: &linebot.EventSource{
					UserID: userId,
				},
				ReplyToken: replyToken,
			},
			expect: nil,
		},
		{
			description: "Invalid Event(Template Message)",
			event: &linebot.Event{
				Message: linebot.NewTemplateMessage("", nil),
				Source: &linebot.EventSource{
					UserID: userId,
				},
				ReplyToken: replyToken,
			},
			expect: nil,
		},
		{
			description: "Invalid Event(Image Map Message)",
			event: &linebot.Event{
				Message: linebot.NewImagemapMessage("", "", linebot.ImagemapBaseSize{Width: 0, Height: 0}, nil),
				Source: &linebot.EventSource{
					UserID: userId,
				},
				ReplyToken: replyToken,
			},
			expect: nil,
		},
		{
			description: "Invalid Event(nil event)",
			event: &linebot.Event{
				Source: &linebot.EventSource{
					UserID: userId,
				},
				ReplyToken: replyToken,
			},
			expect: nil,
		},
	}

	// テスト実行
	for _, test := range tests {
		ret := ec.ConvertEventToEncryptedOperation(test.event)
		if ret != nil && !ret.IsEqual(*test.expect) {
			t.Errorf(testutil.STRING_TEST_MSG_FMT, test.description, test.expect, ret)
		}
		if ret == nil && test.expect != nil {
			t.Errorf(testutil.STRING_TEST_MSG_FMT, test.description, test.expect, "nil")
		}
	}
}
