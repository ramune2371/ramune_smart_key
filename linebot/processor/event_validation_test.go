package processor

import (
	mock_user_info "linebot/dao/user_info/mock"
	"linebot/entity"
	mock_security "linebot/security/mock"
	"testing"
	"time"

	"linebot/testutil"

	"github.com/golang/mock/gomock"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// テスト用のダミー linebot.TextMessage を生成する関数
func createTextMessage(text, userId, replyToken string) *linebot.Event {
	return &linebot.Event{
		Message: &linebot.TextMessage{
			Text: text,
		},
		ReplyToken: replyToken,
		Source: &linebot.EventSource{
			UserID: userId,
		},
	}
}

func TestIsTextMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserInfoDao := mock_user_info.NewMockUserInfoDao(ctrl)
	validator := EventValidatorImpl{UserInfoDao: mockUserInfoDao}

	tests := []struct {
		description string
		event       *linebot.Event
		expect      bool
	}{
		{
			description: "Valid Text Event",
			event: &linebot.Event{
				Message: linebot.NewTextMessage("test"),
			},
			expect: true,
		},
		{
			description: "Invalid Text Event(empty)",
			event: &linebot.Event{
				Message: linebot.NewTextMessage(""),
			},
			expect: false,
		},
		{
			description: "Invalid Event(Audio Message)",
			event: &linebot.Event{
				Message: linebot.NewAudioMessage("test", 0),
			},
			expect: false,
		},
		{
			description: "Invalid Event(Flex Message)",
			event: &linebot.Event{
				Message: linebot.NewFlexMessage("test", nil),
			},
			expect: false,
		},
		{
			description: "Invalid Event(Image Message)",
			event: &linebot.Event{
				Message: linebot.NewImageMessage("test", "test"),
			},
			expect: false,
		},
		{
			description: "Invalid Event(Video Message)",
			event: &linebot.Event{
				Message: linebot.NewVideoMessage("", ""),
			},
			expect: false,
		},
		{
			description: "Invalid Event(Location Message)",
			event: &linebot.Event{
				Message: linebot.NewLocationMessage("title", "", 0, 0),
			},
			expect: false,
		},
		{
			description: "Invalid Event(Sticker Message)",
			event: &linebot.Event{
				Message: linebot.NewStickerMessage("", ""),
			},
			expect: false,
		},
		{
			description: "Invalid Event(Template Message)",
			event: &linebot.Event{
				Message: linebot.NewTemplateMessage("", nil),
			},
			expect: false,
		},
		{
			description: "Invalid Event(Image Map Message)",
			event: &linebot.Event{
				Message: linebot.NewImagemapMessage("", "", linebot.ImagemapBaseSize{Width: 0, Height: 0}, nil),
			},
			expect: false,
		},
		{
			description: "Invalid Event(nil event)",
			event:       &linebot.Event{},
			expect:      false,
		},
	}

	// テスト実行
	for _, test := range tests {
		ret := validator.isTextMessage(test.event)
		if ret != test.expect {
			t.Errorf(testutil.BOOL_TEST_MSG_FMT, test.description, test.expect, ret)
		}
	}
}

func TestVerifyMessageText(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUserInfoDao := mock_user_info.NewMockUserInfoDao(ctrl)
	validator := EventValidatorImpl{UserInfoDao: mockUserInfoDao}

	tests := []struct {
		description string
		text        string
		expect      bool
	}{
		{
			description: "Valid Text(open)",
			text:        "open",
			expect:      true,
		},
		{
			description: "Valid Text(close)",
			text:        "close",
			expect:      true,
		},
		{
			description: "Valid Text(check)",
			text:        "check",
			expect:      true,
		},
		{
			description: "Invalid Text( valid + invalid )",
			text:        "openfizz",
			expect:      false,
		},
		{
			description: "Invalid Text( invalid + valid )",
			text:        "fizzopen",
			expect:      false,
		},
		{
			description: "Invalid Text( invalid + valid + invalid )",
			text:        "fizzopenbuzz",
			expect:      false,
		},
		{
			description: "Invalid Text( valid + space)",
			text:        "open ",
			expect:      false,
		},
		{
			description: "Invalid Text( space + valid )",
			text:        " open",
			expect:      false,
		},
		{
			description: "Invalid Text( space + valid + space )",
			text:        " open ",
			expect:      false,
		},
		{
			description: "Invalid Text(Upper Camel)",
			text:        "Open",
			expect:      false,
		},
		{
			description: "Invalid Text(All Upper)",
			text:        "OPEN",
			expect:      false,
		},
		{
			description: "Invalid Text(empty)",
			text:        "",
			expect:      false,
		},
		{
			description: "Invalid Text(invalid)",
			text:        "fizz",
			expect:      false,
		},
	}

	// テスト実行
	for _, test := range tests {
		ret := validator.verifyMessageText(test.text)
		if ret != test.expect {
			t.Errorf(testutil.BOOL_TEST_MSG_FMT, test.description, test.expect, ret)
		}
	}
}

func TestVerifyUser(t *testing.T) {

	// テスト用のダミー UserInfoDao を作成
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	time := time.Now()
	ValidUser := entity.UserInfo{UserUuid: "validUser", LineId: "valid_user", UserName: "ValidUser", LastAccess: &time, Active: true}
	InvalidUser := entity.UserInfo{UserUuid: "invalidUser", LineId: "invalid_user", UserName: "ValidUser", LastAccess: &time, Active: false}

	mockUserInfoDao := mock_user_info.NewMockUserInfoDao(ctrl)
	mockUserInfoDao.EXPECT().GetUserByLineId("valid_user").Return(&ValidUser)
	mockUserInfoDao.EXPECT().GetUserByLineId("invalid_user").Return(&InvalidUser)
	mockUserInfoDao.EXPECT().GetUserByLineId("null_user").Return(nil)

	validator := EventValidatorImpl{
		UserInfoDao: mockUserInfoDao,
	}

	// テストケースの定義
	tests := []struct {
		description string
		userId      string
		expected    bool
	}{
		{
			description: "Valid(Active) User",
			userId:      "valid_user",
			expected:    true,
		},
		{
			description: "Invalid(Not Active) User",
			userId:      "invalid_user",
			expected:    false,
		},
		{
			description: "Invalid(Not Registered) User",
			userId:      "null_user",
			expected:    false,
		},
	}

	// テスト実行
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			res := validator.verifyUser(test.userId)
			if res != test.expected {
				t.Errorf(testutil.BOOL_TEST_MSG_FMT, test.description, test.expected, res)
			}
		})
	}
}

func TestCheckMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserInfoDao := mock_user_info.NewMockUserInfoDao(ctrl)
	mockEncryptor := mock_security.NewMockEncryptor(ctrl)
	mockEncryptor.EXPECT().SaltHash(gomock.Any()).DoAndReturn(
		func(value string) string {
			return value
		},
	).AnyTimes()
	validator := EventValidatorImpl{UserInfoDao: mockUserInfoDao, Encryptor: mockEncryptor}

	openEvent := createTextMessage("open", "openUser", "openToken")
	openEvent2 := createTextMessage("open", "openUser2", "openToken2")
	closeEvent := createTextMessage("close", "closeUser", "closeToken")
	checkEvent := createTextMessage("check", "checkUser", "checkToken")
	emptyEvent := createTextMessage("", "emptyUser", "emptyToken")
	invalidTextEvent := createTextMessage("fizz", "emptyUser", "emptyToken")
	invalidEvent := &linebot.Event{}

	openOperation := &entity.Operation{OperationId: -1, UserId: "openUser", Operation: entity.Open, ReplyToken: "openToken"}
	openOperation2 := &entity.Operation{OperationId: -1, UserId: "openUser2", Operation: entity.Open, ReplyToken: "openToken2"}
	closeOperation := &entity.Operation{OperationId: -1, UserId: "closeUser", Operation: entity.Close, ReplyToken: "closeToken"}
	checkOperation := &entity.Operation{OperationId: -1, UserId: "checkUser", Operation: entity.Check, ReplyToken: "checkToken"}

	tests := []struct {
		description string
		target      []*linebot.Event
		expect      []*entity.Operation
	}{
		{
			description: "empty",
			target:      nil,
			expect:      nil,
		},
		{
			description: "all ok",
			target: []*linebot.Event{
				openEvent,
				closeEvent,
				checkEvent,
			},
			expect: []*entity.Operation{
				openOperation,
				closeOperation,
				checkOperation,
			},
		},
		{
			description: "include invalid event",
			target: []*linebot.Event{
				openEvent,
				closeEvent,
				checkEvent,
				emptyEvent,
				invalidTextEvent,
				invalidEvent,
			},
			expect: []*entity.Operation{
				openOperation,
				closeOperation,
				checkOperation,
			},
		},
		{
			description: "duplicate",
			target: []*linebot.Event{
				openEvent,
				openEvent,
				openEvent,
				openEvent2,
			},
			expect: []*entity.Operation{
				openOperation,
				openOperation,
				openOperation,
				openOperation2,
			},
		},
	}

	// テスト実行
	for _, test := range tests {
		rets := validator.checkMessage(test.target)
		// length check
		if len(rets) != len(test.expect) {
			t.Errorf("length check:"+testutil.INT_TEST_MSG_FMT, test.description, len(test.expect), len(rets))
		}
		// struct check
		for i, ret := range rets {
			if !ret.IsEqual(*test.expect[i]) {
				t.Errorf("object check:"+testutil.STRING_TEST_MSG_FMT, test.description, *test.expect[i], ret)
			}
		}
	}
}

func TestValidateEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validUser := entity.UserInfo{UserUuid: "valid-uuid", LineId: "valid", UserName: "valid user", LastAccess: nil, Active: true}
	invalidUser := entity.UserInfo{UserUuid: "invalid-uuid", LineId: "invalid", UserName: "invalid user", LastAccess: nil, Active: false}

	mockUserInfoDao := mock_user_info.NewMockUserInfoDao(ctrl)
	mockUserInfoDao.EXPECT().GetUserByLineId("valid").Return(&validUser).AnyTimes()
	mockUserInfoDao.EXPECT().GetUserByLineId("invalid").Return(&invalidUser).AnyTimes()
	mockEncryptor := mock_security.NewMockEncryptor(ctrl)
	mockEncryptor.EXPECT().SaltHash(gomock.Any()).DoAndReturn(func(value string) string { return value }).AnyTimes()
	validator := EventValidatorImpl{UserInfoDao: mockUserInfoDao, Encryptor: mockEncryptor}

	openValidUserEvent := createTextMessage("open", "valid", "openValidToken")
	closeValidUserEvent := createTextMessage("close", "valid", "closeValidToken")
	checkValidUserEvent := createTextMessage("check", "valid", "checkValidToken")

	openInvalidUserEvent := createTextMessage("open", "invalid", "openInvalidToken")
	closeInvalidUserEvent := createTextMessage("close", "invalid", "closeInvalidToken")
	checkInvalidUserEvent := createTextMessage("check", "invalid", "checkInvalidToken")

	emptyEvent := createTextMessage("", "valid", "emptyToken")
	invalidEvent := &linebot.Event{}

	openValidOperation := &entity.Operation{OperationId: -1, UserId: "valid", Operation: entity.Open, ReplyToken: "openValidToken"}
	closeValidOperation := &entity.Operation{OperationId: -1, UserId: "valid", Operation: entity.Close, ReplyToken: "closeValidToken"}
	checkValidOperation := &entity.Operation{OperationId: -1, UserId: "valid", Operation: entity.Check, ReplyToken: "checkValidToken"}

	openInvalidOperation := &entity.Operation{OperationId: -1, UserId: "invalid", Operation: entity.Open, ReplyToken: "openInvalidToken"}
	closeInvalidOperation := &entity.Operation{OperationId: -1, UserId: "invalid", Operation: entity.Close, ReplyToken: "closeInvalidToken"}
	checkInvalidOperation := &entity.Operation{OperationId: -1, UserId: "invalid", Operation: entity.Check, ReplyToken: "checkInvalidToken"}

	target := []*linebot.Event{
		openValidUserEvent,
		closeValidUserEvent,
		checkValidUserEvent,
		openInvalidUserEvent,
		closeInvalidUserEvent,
		checkInvalidUserEvent,
		checkInvalidUserEvent,
		emptyEvent,
		invalidEvent,
	}

	expectValidOperations := []*entity.Operation{
		openValidOperation,
		closeValidOperation,
		checkValidOperation,
	}

	expectInvalidOperations := []*entity.Operation{
		openInvalidOperation,
		closeInvalidOperation,
		checkInvalidOperation,
		checkInvalidOperation,
	}

	// テスト実行
	validRets, invalidRets := validator.ValidateEvent(target)

	// length check
	if len(validRets) != len(expectValidOperations) {
		t.Errorf("valid operation length check:"+testutil.INT_TEST_MSG_FMT, "", len(expectValidOperations), len(validRets))
	}
	if len(invalidRets) != len(expectInvalidOperations) {
		t.Errorf("invalid operation length check:"+testutil.INT_TEST_MSG_FMT, "", len(expectInvalidOperations), len(invalidRets))
	}

	// struct check
	for i, ret := range validRets {
		if !ret.IsEqual(*expectValidOperations[i]) {
			t.Errorf("valid struct check [%d]"+testutil.STRING_TEST_MSG_FMT, i, "", expectValidOperations[i], ret)
		}
	}
	for i, ret := range invalidRets {
		if !ret.IsEqual(*expectInvalidOperations[i]) {
			t.Errorf("invalid struct check [%d]"+testutil.STRING_TEST_MSG_FMT, i, "", expectInvalidOperations[i], ret)
		}
	}
}
