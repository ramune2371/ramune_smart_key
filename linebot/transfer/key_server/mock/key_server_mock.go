// Code generated by MockGen. DO NOT EDIT.
// Source: key_server.go

// Package mock_key_server is a generated GoMock package.
package mock_key_server

import (
	entity "linebot/entity"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockKeyServerTransfer is a mock of KeyServerTransfer interface.
type MockKeyServerTransfer struct {
	ctrl     *gomock.Controller
	recorder *MockKeyServerTransferMockRecorder
}

// MockKeyServerTransferMockRecorder is the mock recorder for MockKeyServerTransfer.
type MockKeyServerTransferMockRecorder struct {
	mock *MockKeyServerTransfer
}

// NewMockKeyServerTransfer creates a new mock instance.
func NewMockKeyServerTransfer(ctrl *gomock.Controller) *MockKeyServerTransfer {
	mock := &MockKeyServerTransfer{ctrl: ctrl}
	mock.recorder = &MockKeyServerTransferMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKeyServerTransfer) EXPECT() *MockKeyServerTransferMockRecorder {
	return m.recorder
}

// CheckKey mocks base method.
func (m *MockKeyServerTransfer) CheckKey() (entity.KeyServerResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckKey")
	ret0, _ := ret[0].(entity.KeyServerResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckKey indicates an expected call of CheckKey.
func (mr *MockKeyServerTransferMockRecorder) CheckKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckKey", reflect.TypeOf((*MockKeyServerTransfer)(nil).CheckKey))
}

// CloseKey mocks base method.
func (m *MockKeyServerTransfer) CloseKey() (entity.KeyServerResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseKey")
	ret0, _ := ret[0].(entity.KeyServerResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloseKey indicates an expected call of CloseKey.
func (mr *MockKeyServerTransferMockRecorder) CloseKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseKey", reflect.TypeOf((*MockKeyServerTransfer)(nil).CloseKey))
}

// OpenKey mocks base method.
func (m *MockKeyServerTransfer) OpenKey() (entity.KeyServerResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenKey")
	ret0, _ := ret[0].(entity.KeyServerResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenKey indicates an expected call of OpenKey.
func (mr *MockKeyServerTransferMockRecorder) OpenKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenKey", reflect.TypeOf((*MockKeyServerTransfer)(nil).OpenKey))
}
