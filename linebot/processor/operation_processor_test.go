package processor

import (
	"fmt"
	"linebot/applicationerror"
	"linebot/dao/operation_history"
	mock_operation_history "linebot/dao/operation_history/mock"
	"linebot/dao/user_info"
	mock_user_info "linebot/dao/user_info/mock"
	"linebot/entity"
	"linebot/entity/message"
	"linebot/testutil"
	"linebot/transfer/key_server"
	mock_key_server "linebot/transfer/key_server/mock"
	"linebot/transfer/line"
	mock_line "linebot/transfer/line/mock"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type MockEncryptor struct{}

func (MockEncryptor) SaltHash(value string) string { return value + "SaltHash" }

func getFailMockedKeyServerTransfer(t *testing.T, ctrl *gomock.Controller) key_server.KeyServerTransfer {
	// failMockedKeyServerTransfer
	failFunc := func() {
		t.Fatalf("Not expected call!!")
	}
	failMockedKeyServerTransfer := mock_key_server.NewMockKeyServerTransfer(ctrl)
	failMockedKeyServerTransfer.EXPECT().CloseKey().DoAndReturn(failFunc).AnyTimes()
	failMockedKeyServerTransfer.EXPECT().OpenKey().DoAndReturn(failFunc).AnyTimes()
	failMockedKeyServerTransfer.EXPECT().CheckKey().DoAndReturn(failFunc).AnyTimes()

	return failMockedKeyServerTransfer
}

func createMockKeyServerTransfer(operationType string, t *testing.T, ctrl *gomock.Controller, res entity.KeyServerResponse, err *applicationerror.ApplicationError) key_server.KeyServerTransfer {
	helperFunc := func(expectOp string) (entity.KeyServerResponse, error) {
		if operationType != expectOp {
			t.Fatalf("Not expected Called")
		}

		if err == nil {
			return res, nil
		}
		return res, err
	}
	// key server transfer のmock化
	mockedKeyServerTransfer := mock_key_server.NewMockKeyServerTransfer(ctrl)
	mockedKeyServerTransfer.EXPECT().CloseKey().DoAndReturn(func() (entity.KeyServerResponse, error) {
		return helperFunc("close")
	}).AnyTimes()
	mockedKeyServerTransfer.EXPECT().OpenKey().DoAndReturn(func() (entity.KeyServerResponse, error) {
		return helperFunc("open")
	}).AnyTimes()
	mockedKeyServerTransfer.EXPECT().CheckKey().DoAndReturn(func() (entity.KeyServerResponse, error) {
		return helperFunc("check")
	}).AnyTimes()

	return mockedKeyServerTransfer
}

func createSingleAssertLineTransfer(assert func(resText, replyToken string) error, t *testing.T, ctrl *gomock.Controller) line.LineTransfer {
	// line transfer のmock化
	mockedLineTransfer := mock_line.NewMockLineTransfer(ctrl)
	mockedLineTransfer.EXPECT().ReplyToToken(gomock.Any(), gomock.Any()).DoAndReturn(assert).AnyTimes()
	return mockedLineTransfer
}

func getMockedUserInfoDao(t *testing.T, ctrl *gomock.Controller) user_info.UserInfoDao {
	operationTypes := []string{"open", "close", "check"}
	mockedUserInfoDao := mock_user_info.NewMockUserInfoDao(ctrl)
	mockedUserInfoDao.EXPECT().GetUserByLineId(gomock.Any()).DoAndReturn(
		func(lineId string) *entity.UserInfo {
			return &entity.UserInfo{
				UserUuid:   "",
				LineId:     lineId,
				UserName:   lineId,
				LastAccess: nil,
				Active:     !strings.Contains(lineId, "invalid"),
			}
		},
	).AnyTimes()
	for _, o := range operationTypes {
		mockedUserInfoDao.EXPECT().UpsertInvalidUser("invalid" + o + "SaltHash").Return(true).AnyTimes()
	}
	for _, o := range operationTypes {
		mockedUserInfoDao.EXPECT().UpdateUserLastAccess("valid" + o + "SaltHash").Return(true).AnyTimes()
	}
	return mockedUserInfoDao
}

func getMockOperationHistoryDao(t *testing.T, ctrl *gomock.Controller) operation_history.OperationHistoryDao {
	mockedOperationHistoryDao := mock_operation_history.NewMockOperationHistoryDao(ctrl)
	id := -1
	mockedOperationHistoryDao.EXPECT().InsertOperationHistory(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(lineId string, operationType entity.OperationType, operationResult entity.OperationResult) *entity.OperationHistory {
		return &entity.OperationHistory{
			OperationId:     &id,
			LineId:          lineId,
			OperationType:   operationType,
			OperationResult: operationResult,
			OperationTime:   nil,
		}
	}).AnyTimes()
	mockedOperationHistoryDao.EXPECT().UpdateOperationHistoryByOperationId(gomock.Any(), gomock.Any()).Return(-1).AnyTimes()
	mockedOperationHistoryDao.EXPECT().UpdateOperationHistoryWithErrorByOperationId(gomock.Any(), gomock.Any()).Return(-1).AnyTimes()
	return mockedOperationHistoryDao
}

func TestOperationProcessor_HandleEvents(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	operationTypes := []string{"open", "close", "check"}

	t.Run("正常系(isOperating:false, 単一ユーザ)", func(t *testing.T) {

		t.Run("有効ユーザのみからのリクエスト", func(t *testing.T) {
			expectResTexts := map[string]string{
				"open":  message.SUCCESS_OPEN,
				"close": message.SUCCESS_CLOSE,
				"check": message.CHECK_OPEN,
			}
			response := map[string]entity.KeyServerResponse{
				"open":  {KeyStatus: entity.KeyStatusOpen, OperationStatus: entity.OperationComplete},
				"close": {KeyStatus: entity.KeyStatusClose, OperationStatus: entity.OperationComplete},
				"check": {KeyStatus: entity.KeyStatusOpen, OperationStatus: entity.OperationComplete},
			}
			for _, ot := range operationTypes {
				targetEvents := []*linebot.Event{
					{Message: linebot.NewTextMessage(ot), Source: &linebot.EventSource{UserID: "valid" + ot}, ReplyToken: "validToken"},
				}

				// key server transfer のmock化
				mockedKeyServerTransfer := createMockKeyServerTransfer(ot, t, ctrl, response[ot], nil)

				// line transfer のmock化
				assert := func(resText, replyToken string) error {
					if replyToken == "validToken" {
						if resText != expectResTexts[ot] {
							t.Errorf(testutil.STRING_TEST_MSG_FMT, "", expectResTexts[ot], resText)
						}
					} else {
						t.Errorf("ユーザ無効判定(想定外Assert)")
					}
					return nil
				}
				mockedLineTransfer := createSingleAssertLineTransfer(assert, t, ctrl)

				opProcessor := NewOperationProcessor(
					getMockOperationHistoryDao(t, ctrl),
					getMockedUserInfoDao(t, ctrl),
					mockedLineTransfer,
					mockedKeyServerTransfer,
					MockEncryptor{},
				)
				// HandleEventsを呼び出す
				opProcessor.SetIsOperating(false)
				opProcessor.HandleEvents(targetEvents)
			}
		})

		t.Run("無効ユーザのみからのリクエスト", func(t *testing.T) {
			for _, ot := range operationTypes {
				targetEvents := []*linebot.Event{
					{Message: linebot.NewTextMessage(ot), Source: &linebot.EventSource{UserID: "invalid" + ot}, ReplyToken: "invalidToken"},
				}

				// line transfer のmock化
				assert := func(resText, replyToken string) error {
					if replyToken == "invalidToken" {
						expectText := "無効なユーザだよ。↓の文字列を管理者に送って。\n「invalid" + ot + "SaltHash」"
						if resText != expectText {
							t.Errorf(testutil.STRING_TEST_MSG_FMT, "", expectText, resText)
						}
					} else {
						t.Errorf("ユーザ有効判定(想定外Assert)")
					}
					return nil
				}
				mockedLineTransfer := createSingleAssertLineTransfer(assert, t, ctrl)

				opProcessor := NewOperationProcessor(
					getMockOperationHistoryDao(t, ctrl),
					getMockedUserInfoDao(t, ctrl),
					mockedLineTransfer,
					getFailMockedKeyServerTransfer(t, ctrl),
					MockEncryptor{},
				)
				// HandleEventsを呼び出す
				opProcessor.SetIsOperating(false)
				opProcessor.HandleEvents(targetEvents)
			}
		})

	})
	t.Run("正常系(isOperating:true、単一ユーザ)", func(t *testing.T) {

		t.Run("有効ユーザのみからのリクエスト", func(t *testing.T) {
			expectResTexts := message.ANOTHER_OPERATING
			for _, ot := range operationTypes {
				targetEvents := []*linebot.Event{
					{Message: linebot.NewTextMessage(ot), Source: &linebot.EventSource{UserID: "valid" + ot}, ReplyToken: "validToken"},
				}

				// line transfer のmock化
				assert := func(resText, replyToken string) error {
					if replyToken == "validToken" {
						if resText != expectResTexts {
							t.Errorf(testutil.STRING_TEST_MSG_FMT, "", expectResTexts, resText)
						}
					} else {
						t.Errorf("ユーザ無効判定(想定外Assert)")
					}
					return nil
				}
				mockedLineTransfer := createSingleAssertLineTransfer(assert, t, ctrl)

				opProcessor := NewOperationProcessor(
					getMockOperationHistoryDao(t, ctrl),
					getMockedUserInfoDao(t, ctrl),
					mockedLineTransfer,
					getFailMockedKeyServerTransfer(t, ctrl),
					MockEncryptor{},
				)
				// HandleEventsを呼び出す
				opProcessor.SetIsOperating(true)
				opProcessor.HandleEvents(targetEvents)
			}
		})

		t.Run("無効ユーザのみからのリクエスト", func(t *testing.T) {
			for _, ot := range operationTypes {
				expectResTexts := "無効なユーザだよ。↓の文字列を管理者に送って。\n「invalid" + ot + "SaltHash」"
				targetEvents := []*linebot.Event{
					{Message: linebot.NewTextMessage(ot), Source: &linebot.EventSource{UserID: "invalid" + ot}, ReplyToken: "invalidToken"},
				}

				// line transfer のmock化
				assert := func(resText, replyToken string) error {
					if replyToken == "invalidToken" {
						if resText != expectResTexts {
							t.Errorf(testutil.STRING_TEST_MSG_FMT, "", expectResTexts, resText)
						}
					} else {
						t.Errorf("ユーザ有効判定(想定外Assert)")
					}
					return nil
				}
				mockedLineTransfer := createSingleAssertLineTransfer(assert, t, ctrl)

				opProcessor := NewOperationProcessor(
					getMockOperationHistoryDao(t, ctrl),
					getMockedUserInfoDao(t, ctrl),
					mockedLineTransfer,
					getFailMockedKeyServerTransfer(t, ctrl),
					MockEncryptor{},
				)
				// HandleEventsを呼び出す
				opProcessor.SetIsOperating(true)
				opProcessor.HandleEvents(targetEvents)
			}
		})
	})

	t.Run("異常系(鍵サーバとの通信異常系)", func(t *testing.T) {
		errors := []*applicationerror.ApplicationError{
			applicationerror.ConnectionError,
			applicationerror.ResponseParseError,
		}
		expectResTexts := []string{
			message.CONNECTION_ERROR,
			fmt.Sprintf(message.APPLICATION_ERROR, "102"),
		}
		for i, err := range errors {
			expectResText := expectResTexts[i]
			for _, ot := range operationTypes {
				targetEvents := []*linebot.Event{
					{Message: linebot.NewTextMessage(ot), Source: &linebot.EventSource{UserID: "valid" + ot}, ReplyToken: "validToken"},
				}

				// key server transfer のmock化
				mockedKeyServerTransfer := createMockKeyServerTransfer(ot, t, ctrl, entity.KeyServerResponse{}, err)

				// line transfer のmock化
				mockedLineTransfer := mock_line.NewMockLineTransfer(ctrl)
				mockedLineTransfer.EXPECT().ReplyToToken(gomock.Any(), gomock.Any()).DoAndReturn(func(resText, replyToken string) error {
					if replyToken == "validToken" {
						if resText != expectResText {
							t.Errorf(testutil.STRING_TEST_MSG_FMT, "", expectResText, resText)
						}
					} else {
						t.Errorf("ユーザ無効判定(想定外Assert)")
					}
					return nil
				}).AnyTimes()
				opProcessor := NewOperationProcessor(
					getMockOperationHistoryDao(t, ctrl),
					getMockedUserInfoDao(t, ctrl),
					mockedLineTransfer,
					mockedKeyServerTransfer,
					MockEncryptor{},
				)
				// HandleEventsを呼び出す
				opProcessor.SetIsOperating(false)
				opProcessor.HandleEvents(targetEvents)
			}
		}
	})
}

func TestHandleMergeEvents(t *testing.T) {
	user1Open := entity.Operation{OperationId: -1, UserId: "user1", Operation: entity.Open, ReplyToken: ""}
	user1Close := entity.Operation{OperationId: -1, UserId: "user1", Operation: entity.Close, ReplyToken: ""}
	user1Check := entity.Operation{OperationId: -1, UserId: "user1", Operation: entity.Check, ReplyToken: ""}
	user2Open := entity.Operation{OperationId: -1, UserId: "user2", Operation: entity.Open, ReplyToken: ""}
	user2Close := entity.Operation{OperationId: -1, UserId: "user2", Operation: entity.Close, ReplyToken: ""}
	user2Check := entity.Operation{OperationId: -1, UserId: "user2", Operation: entity.Check, ReplyToken: ""}

	tests := []struct {
		description           string
		targetEvents          []*entity.Operation
		expectOperations      map[string]entity.Operation
		expectMasterOperation entity.OperationType
	}{
		{
			description:           "Event数1(open)",
			targetEvents:          []*entity.Operation{&user1Open},
			expectOperations:      map[string]entity.Operation{"user1": user1Open},
			expectMasterOperation: entity.Open,
		},
		{
			description:           "Event数1(close)",
			targetEvents:          []*entity.Operation{&user1Close},
			expectOperations:      map[string]entity.Operation{"user1": user1Close},
			expectMasterOperation: entity.Close,
		},
		{
			description:           "Event数1(check)",
			targetEvents:          []*entity.Operation{&user1Check},
			expectOperations:      map[string]entity.Operation{"user1": user1Check},
			expectMasterOperation: entity.Check,
		},
		{
			description:           "Event数2(複数ユーザ(open,open))",
			targetEvents:          []*entity.Operation{&user1Open, &user2Open},
			expectOperations:      map[string]entity.Operation{"user1": user1Open, "user2": user2Open},
			expectMasterOperation: entity.Open,
		},
		{
			description:           "Event数2(複数ユーザ(open,close))",
			targetEvents:          []*entity.Operation{&user1Open, &user2Close},
			expectOperations:      map[string]entity.Operation{"user1": user1Open, "user2": user2Close},
			expectMasterOperation: entity.Close,
		},
		{
			description:           "Event数2(複数ユーザ(open,check))",
			targetEvents:          []*entity.Operation{&user1Open, &user2Check},
			expectOperations:      map[string]entity.Operation{"user1": user1Open, "user2": user2Check},
			expectMasterOperation: entity.Open,
		},
		{
			description:           "Event数2(複数ユーザ(close,open))",
			targetEvents:          []*entity.Operation{&user1Close, &user2Open},
			expectOperations:      map[string]entity.Operation{"user1": user1Close, "user2": user2Open},
			expectMasterOperation: entity.Open,
		},
		{
			description:           "Event数2(複数ユーザ(close,close))",
			targetEvents:          []*entity.Operation{&user1Close, &user2Close},
			expectOperations:      map[string]entity.Operation{"user1": user1Close, "user2": user2Close},
			expectMasterOperation: entity.Close,
		},
		{
			description:           "Event数2(複数ユーザ(close,check))",
			targetEvents:          []*entity.Operation{&user1Close, &user2Check},
			expectOperations:      map[string]entity.Operation{"user1": user1Close, "user2": user2Check},
			expectMasterOperation: entity.Close,
		},
		{
			description:           "Event数2(複数ユーザ(check,open))",
			targetEvents:          []*entity.Operation{&user1Check, &user2Open},
			expectOperations:      map[string]entity.Operation{"user1": user1Check, "user2": user2Open},
			expectMasterOperation: entity.Open,
		},
		{
			description:           "Event数2(複数ユーザ(check,close))",
			targetEvents:          []*entity.Operation{&user1Check, &user2Close},
			expectOperations:      map[string]entity.Operation{"user1": user1Check, "user2": user2Close},
			expectMasterOperation: entity.Close,
		},
		{
			description:           "Event数2(複数ユーザ(check,check))",
			targetEvents:          []*entity.Operation{&user1Check, &user2Check},
			expectOperations:      map[string]entity.Operation{"user1": user1Check, "user2": user2Check},
			expectMasterOperation: entity.Check,
		},
		{
			description:           "Event数2(同一ユーザ(open,open))",
			targetEvents:          []*entity.Operation{&user1Open, &user1Open},
			expectOperations:      map[string]entity.Operation{"user1": user1Open},
			expectMasterOperation: entity.Open,
		},
		{
			description:           "Event数2(同一ユーザ(open,close))",
			targetEvents:          []*entity.Operation{&user1Open, &user1Close},
			expectOperations:      map[string]entity.Operation{"user1": user1Close},
			expectMasterOperation: entity.Close,
		},
		{
			description:           "Event数2(同一ユーザ(open,check))",
			targetEvents:          []*entity.Operation{&user1Open, &user1Check},
			expectOperations:      map[string]entity.Operation{"user1": user1Open},
			expectMasterOperation: entity.Open,
		},
		{
			description:           "Event数2(同一ユーザ(close,open))",
			targetEvents:          []*entity.Operation{&user1Close, &user1Open},
			expectOperations:      map[string]entity.Operation{"user1": user1Open},
			expectMasterOperation: entity.Open,
		},
		{
			description:           "Event数2(同一ユーザ(close,close))",
			targetEvents:          []*entity.Operation{&user1Close, &user1Close},
			expectOperations:      map[string]entity.Operation{"user1": user1Close},
			expectMasterOperation: entity.Close,
		},
		{
			description:           "Event数2(同一ユーザ(close,check))",
			targetEvents:          []*entity.Operation{&user1Close, &user1Check},
			expectOperations:      map[string]entity.Operation{"user1": user1Close},
			expectMasterOperation: entity.Close,
		},
		{
			description:           "Event数2(同一ユーザ(check,open))",
			targetEvents:          []*entity.Operation{&user1Check, &user1Open},
			expectOperations:      map[string]entity.Operation{"user1": user1Open},
			expectMasterOperation: entity.Open,
		},
		{
			description:           "Event数2(同一ユーザ(check,close))",
			targetEvents:          []*entity.Operation{&user1Check, &user1Close},
			expectOperations:      map[string]entity.Operation{"user1": user1Close},
			expectMasterOperation: entity.Close,
		},
		{
			description:           "Event数2(同一ユーザ(check,check))",
			targetEvents:          []*entity.Operation{&user1Check, &user1Check},
			expectOperations:      map[string]entity.Operation{"user1": user1Check},
			expectMasterOperation: entity.Check,
		},
		{
			description:           "Event数3(同一ユーザ(open,check,close))<過剰テスト>",
			targetEvents:          []*entity.Operation{&user1Open, &user1Check, &user1Close},
			expectOperations:      map[string]entity.Operation{"user1": user1Close},
			expectMasterOperation: entity.Close,
		},
		{
			description:           "Event数3(複数ユーザ(open,check,close))<過剰テスト>",
			targetEvents:          []*entity.Operation{&user1Open, &user2Check, &user1Close},
			expectOperations:      map[string]entity.Operation{"user1": user1Close, "user2": user2Check},
			expectMasterOperation: entity.Close,
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			op := NewOperationProcessor(
				getMockOperationHistoryDao(t, ctrl),
				getMockedUserInfoDao(t, ctrl),
				createSingleAssertLineTransfer(func(resText, replyToken string) error { return nil }, t, ctrl),
				createMockKeyServerTransfer("open", t, ctrl, entity.KeyServerResponse{}, nil),
				MockEncryptor{},
			)
			ops, masterOp := op.mergeEvents(test.targetEvents)

			for k, o := range ops {
				if test.expectOperations[k] != o {
					t.Errorf(testutil.STRING_TEST_MSG_FMT, test.description, test.expectOperations[k], o)
				}
			}
			if test.expectMasterOperation != masterOp {
				t.Errorf(testutil.INT_TEST_MSG_FMT, test.description, test.expectMasterOperation, masterOp)
			}
		})
	}
}

func TestHandleMasterOperation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		description string
		target      entity.OperationType
	}{
		{
			description: "open",
			target:      entity.Open,
		},
		{
			description: "close",
			target:      entity.Close,
		},
		{
			description: "check",
			target:      entity.Check,
		},
		{
			description: "unsupported",
			target:      entity.Unsupported,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ksTransfer := createMockKeyServerTransfer(test.description, t, ctrl, entity.KeyServerResponse{}, nil)
			op := NewOperationProcessor(
				operation_history.NewEmptyOperationHistoryDao(),
				user_info.NewEmptyUserInfoDao(),
				mock_line.NewMockLineTransfer(ctrl),
				ksTransfer,
				MockEncryptor{},
			)
			op.handleMasterOperation(test.target)
		})
	}
}

func TestReplyCheckResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		description  string
		targetToken  string
		targetResult entity.KeyStatus
		expectText   string
	}{
		{
			description:  "open",
			targetToken:  "openToken",
			targetResult: entity.KeyStatusOpen,
			expectText:   message.CHECK_OPEN,
		},
		{
			description:  "close",
			targetToken:  "closeToken",
			targetResult: entity.KeyStatusClose,
			expectText:   message.CHECK_CLOSE,
		},
	}

	for _, test := range tests {
		op := NewOperationProcessor(
			operation_history.NewEmptyOperationHistoryDao(),
			user_info.NewEmptyUserInfoDao(),
			createSingleAssertLineTransfer(func(resText, replyToken string) error {
				if resText != test.expectText {
					t.Errorf(testutil.STRING_TEST_MSG_FMT, test.description, test.expectText, resText)
				}
				return nil
			}, t, ctrl),
			getFailMockedKeyServerTransfer(t, ctrl),
			MockEncryptor{},
		)
		op.replyCheckResult(test.targetToken, test.targetResult)
	}
}

func TestHandleKeyServerResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user1Open := entity.Operation{OperationId: -1, UserId: "user1", Operation: entity.Open, ReplyToken: "user1Token"}
	user1Close := entity.Operation{OperationId: -1, UserId: "user1", Operation: entity.Close, ReplyToken: "user1Token"}
	user1Check := entity.Operation{OperationId: -1, UserId: "user1", Operation: entity.Check, ReplyToken: "user1Token"}

	tests := []struct {
		description           string
		targetOps             map[string]entity.Operation
		targetKeyStatus       entity.KeyStatus
		targetOperationStatus entity.OperationStatus
		expectResText         string
	}{
		// key status open
		// user1open
		{
			description:           "単一ユーザ(操作:open, 鍵状態:Open, 操作状態:変更なし)",
			targetOps:             map[string]entity.Operation{"user1": user1Open},
			targetKeyStatus:       entity.KeyStatusOpen,
			targetOperationStatus: entity.OperationAlready,
			expectResText:         message.SUCCESS_OPEN,
		},
		{
			description:           "単一ユーザ(操作:open, 鍵状態:Open, 操作状態:成功)",
			targetOps:             map[string]entity.Operation{"user1": user1Open},
			targetKeyStatus:       entity.KeyStatusOpen,
			targetOperationStatus: entity.OperationComplete,
			expectResText:         message.SUCCESS_OPEN,
		},
		{
			description:           "単一ユーザ(操作:open, 鍵状態:Open, 操作状態:他の人が操作中)",
			targetOps:             map[string]entity.Operation{"user1": user1Open},
			targetKeyStatus:       entity.KeyStatusOpen,
			targetOperationStatus: entity.OperationAnother,
			expectResText:         message.ANOTHER_OPERATING,
		},
		{
			description:           "単一ユーザ(操作:open, 鍵状態:Open, 操作状態:unknown)",
			targetOps:             map[string]entity.Operation{"user1": user1Open},
			targetKeyStatus:       entity.KeyStatusOpen,
			targetOperationStatus: entity.OperationUnknown,
			expectResText:         message.SUCCESS_OPEN,
		},
		// user1close
		{
			description:           "単一ユーザ(操作:close, 鍵状態:Open, 操作状態:変更なし)",
			targetOps:             map[string]entity.Operation{"user1": user1Close},
			targetKeyStatus:       entity.KeyStatusOpen,
			targetOperationStatus: entity.OperationAlready,
			expectResText:         message.ANOTHER_OPEN,
		},
		{
			description:           "単一ユーザ(操作:close, 鍵状態:Open, 操作状態:成功)",
			targetOps:             map[string]entity.Operation{"user1": user1Close},
			targetKeyStatus:       entity.KeyStatusOpen,
			targetOperationStatus: entity.OperationComplete,
			expectResText:         message.ANOTHER_OPEN,
		},
		{
			description:           "単一ユーザ(操作:close, 鍵状態:Open, 操作状態:他の人が操作中)",
			targetOps:             map[string]entity.Operation{"user1": user1Close},
			targetKeyStatus:       entity.KeyStatusOpen,
			targetOperationStatus: entity.OperationAnother,
			expectResText:         message.ANOTHER_OPERATING,
		},
		{
			description:           "単一ユーザ(操作:close, 鍵状態:Open, 操作状態:unknown)",
			targetOps:             map[string]entity.Operation{"user1": user1Close},
			targetKeyStatus:       entity.KeyStatusOpen,
			targetOperationStatus: entity.OperationUnknown,
			expectResText:         message.ANOTHER_OPEN,
		},
		// user1check
		{
			description:           "単一ユーザ(操作:check, 鍵状態:Open, 操作状態:変更なし)",
			targetOps:             map[string]entity.Operation{"user1": user1Check},
			targetKeyStatus:       entity.KeyStatusOpen,
			targetOperationStatus: entity.OperationAlready,
			expectResText:         message.CHECK_OPEN,
		},
		{
			description:           "単一ユーザ(操作:check, 鍵状態:Open, 操作状態:成功)",
			targetOps:             map[string]entity.Operation{"user1": user1Check},
			targetKeyStatus:       entity.KeyStatusOpen,
			targetOperationStatus: entity.OperationComplete,
			expectResText:         message.CHECK_OPEN,
		},
		{
			description:           "単一ユーザ(操作:check, 鍵状態:Open, 操作状態:他の人が操作中)",
			targetOps:             map[string]entity.Operation{"user1": user1Check},
			targetKeyStatus:       entity.KeyStatusOpen,
			targetOperationStatus: entity.OperationAnother,
			expectResText:         message.ANOTHER_OPERATING,
		},
		{
			description:           "単一ユーザ(操作:check, 鍵状態:Open, 操作状態:unknown)",
			targetOps:             map[string]entity.Operation{"user1": user1Check},
			targetKeyStatus:       entity.KeyStatusOpen,
			targetOperationStatus: entity.OperationUnknown,
			expectResText:         message.CHECK_OPEN,
		},

		// key status close
		// user1open
		{
			description:           "単一ユーザ(操作:open, 鍵状態:close, 操作状態:変更なし)",
			targetOps:             map[string]entity.Operation{"user1": user1Open},
			targetKeyStatus:       entity.KeyStatusClose,
			targetOperationStatus: entity.OperationAlready,
			expectResText:         message.ANOTHER_CLOSE,
		},
		{
			description:           "単一ユーザ(操作:open, 鍵状態:Open, 操作状態:成功)",
			targetOps:             map[string]entity.Operation{"user1": user1Open},
			targetKeyStatus:       entity.KeyStatusClose,
			targetOperationStatus: entity.OperationComplete,
			expectResText:         message.ANOTHER_CLOSE,
		},
		{
			description:           "単一ユーザ(操作:open, 鍵状態:Open, 操作状態:他の人が操作中)",
			targetOps:             map[string]entity.Operation{"user1": user1Open},
			targetKeyStatus:       entity.KeyStatusClose,
			targetOperationStatus: entity.OperationAnother,
			expectResText:         message.ANOTHER_OPERATING,
		},
		{
			description:           "単一ユーザ(操作:open, 鍵状態:Open, 操作状態:unknown)",
			targetOps:             map[string]entity.Operation{"user1": user1Open},
			targetKeyStatus:       entity.KeyStatusClose,
			targetOperationStatus: entity.OperationUnknown,
			expectResText:         message.ANOTHER_CLOSE,
		},
		// user1close
		{
			description:           "単一ユーザ(操作:close, 鍵状態:Open, 操作状態:変更なし)",
			targetOps:             map[string]entity.Operation{"user1": user1Close},
			targetKeyStatus:       entity.KeyStatusClose,
			targetOperationStatus: entity.OperationAlready,
			expectResText:         message.SUCCESS_CLOSE,
		},
		{
			description:           "単一ユーザ(操作:close, 鍵状態:Open, 操作状態:成功)",
			targetOps:             map[string]entity.Operation{"user1": user1Close},
			targetKeyStatus:       entity.KeyStatusClose,
			targetOperationStatus: entity.OperationComplete,
			expectResText:         message.SUCCESS_CLOSE,
		},
		{
			description:           "単一ユーザ(操作:close, 鍵状態:Open, 操作状態:他の人が操作中)",
			targetOps:             map[string]entity.Operation{"user1": user1Close},
			targetKeyStatus:       entity.KeyStatusClose,
			targetOperationStatus: entity.OperationAnother,
			expectResText:         message.ANOTHER_OPERATING,
		},
		{
			description:           "単一ユーザ(操作:close, 鍵状態:Open, 操作状態:unknown)",
			targetOps:             map[string]entity.Operation{"user1": user1Close},
			targetKeyStatus:       entity.KeyStatusClose,
			targetOperationStatus: entity.OperationUnknown,
			expectResText:         message.SUCCESS_CLOSE,
		},
		// user1check
		{
			description:           "単一ユーザ(操作:check, 鍵状態:Open, 操作状態:変更なし)",
			targetOps:             map[string]entity.Operation{"user1": user1Check},
			targetKeyStatus:       entity.KeyStatusClose,
			targetOperationStatus: entity.OperationAlready,
			expectResText:         message.CHECK_CLOSE,
		},
		{
			description:           "単一ユーザ(操作:check, 鍵状態:Open, 操作状態:成功)",
			targetOps:             map[string]entity.Operation{"user1": user1Check},
			targetKeyStatus:       entity.KeyStatusClose,
			targetOperationStatus: entity.OperationComplete,
			expectResText:         message.CHECK_CLOSE,
		},
		{
			description:           "単一ユーザ(操作:check, 鍵状態:Open, 操作状態:他の人が操作中)",
			targetOps:             map[string]entity.Operation{"user1": user1Check},
			targetKeyStatus:       entity.KeyStatusClose,
			targetOperationStatus: entity.OperationAnother,
			expectResText:         message.ANOTHER_OPERATING,
		},
		{
			description:           "単一ユーザ(操作:check, 鍵状態:Open, 操作状態:unknown)",
			targetOps:             map[string]entity.Operation{"user1": user1Check},
			targetKeyStatus:       entity.KeyStatusClose,
			targetOperationStatus: entity.OperationUnknown,
			expectResText:         message.CHECK_CLOSE,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			op := NewOperationProcessor(
				operation_history.NewEmptyOperationHistoryDao(),
				user_info.NewEmptyUserInfoDao(),
				createSingleAssertLineTransfer(func(resText, replyToken string) error {
					if resText != test.expectResText {
						t.Errorf(testutil.STRING_TEST_MSG_FMT, test.description, test.expectResText, resText)
					}
					return nil
				}, t, ctrl),
				getFailMockedKeyServerTransfer(t, ctrl),
				MockEncryptor{},
			)
			op.handleKeyServerResult(test.targetOps, entity.KeyServerResponse{KeyStatus: test.targetKeyStatus, OperationStatus: test.targetOperationStatus})
		})
	}
}
