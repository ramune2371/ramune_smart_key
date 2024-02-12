// Code generated by MockGen. DO NOT EDIT.
// Source: line.go

// Package mock_line is a generated GoMock package.
package mock_line

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockLineTransfer is a mock of LineTransfer interface.
type MockLineTransfer struct {
	ctrl     *gomock.Controller
	recorder *MockLineTransferMockRecorder
}

// MockLineTransferMockRecorder is the mock recorder for MockLineTransfer.
type MockLineTransferMockRecorder struct {
	mock *MockLineTransfer
}

// NewMockLineTransfer creates a new mock instance.
func NewMockLineTransfer(ctrl *gomock.Controller) *MockLineTransfer {
	mock := &MockLineTransfer{ctrl: ctrl}
	mock.recorder = &MockLineTransferMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLineTransfer) EXPECT() *MockLineTransferMockRecorder {
	return m.recorder
}

// ReplyToToken mocks base method.
func (m *MockLineTransfer) ReplyToToken(resText, replyToken string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplyToToken", resText, replyToken)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReplyToToken indicates an expected call of ReplyToToken.
func (mr *MockLineTransferMockRecorder) ReplyToToken(resText, replyToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplyToToken", reflect.TypeOf((*MockLineTransfer)(nil).ReplyToToken), resText, replyToken)
}