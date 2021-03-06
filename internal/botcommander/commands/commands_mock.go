// Code generated by MockGen. DO NOT EDIT.
// Source: commands.go

// Package mock_commands is a generated GoMock package.
package commands

import (
	"github.com/unanet/eve-bot/internal/chatservice/chatmodels"
	reflect "reflect"

	dynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	gomock "github.com/golang/mock/gomock"
)

// MockEvebotCommand is a mock of EvebotCommand interface.
type MockEvebotCommand struct {
	ctrl     *gomock.Controller
	recorder *MockEvebotCommandMockRecorder
}

// MockEvebotCommandMockRecorder is the mock recorder for MockEvebotCommand.
type MockEvebotCommandMockRecorder struct {
	mock *MockEvebotCommand
}

// NewMockEvebotCommand creates a new mock instance.
func NewMockEvebotCommand(ctrl *gomock.Controller) *MockEvebotCommand {
	mock := &MockEvebotCommand{ctrl: ctrl}
	mock.recorder = &MockEvebotCommandMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEvebotCommand) EXPECT() *MockEvebotCommandMockRecorder {
	return m.recorder
}

// AckMsg mocks base method.
func (m *MockEvebotCommand) AckMsg() (string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AckMsg")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// AckMsg indicates an expected call of AckMsg.
func (mr *MockEvebotCommandMockRecorder) AckMsg() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AckMsg", reflect.TypeOf((*MockEvebotCommand)(nil).AckMsg))
}

// Info mocks base method.
func (m *MockEvebotCommand) Info() ChatInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Info")
	ret0, _ := ret[0].(ChatInfo)
	return ret0
}

// Info indicates an expected call of Info.
func (mr *MockEvebotCommandMockRecorder) Info() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockEvebotCommand)(nil).Info))
}

// IsAuthenticated mocks base method.
func (m *MockEvebotCommand) IsAuthenticated(chatUserFn *chatmodels.ChatUser, db *dynamodb.DynamoDB) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsAuthenticated", chatUserFn, db)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsAuthenticated indicates an expected call of IsAuthenticated.
func (mr *MockEvebotCommandMockRecorder) IsAuthenticated(chatUserFn, db interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsAuthenticated", reflect.TypeOf((*MockEvebotCommand)(nil).IsAuthenticated), chatUserFn, db)
}

// IsAuthorized mocks base method.
func (m *MockEvebotCommand) IsAuthorized(allowedChannel map[string]interface{}, chatChanFn ChatChannelInfoFn) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "isAuthorized", allowedChannel, chatChanFn)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsAuthorized indicates an expected call of IsAuthorized.
func (mr *MockEvebotCommandMockRecorder) IsAuthorized(allowedChannel, chatChanFn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "isAuthorized", reflect.TypeOf((*MockEvebotCommand)(nil).IsAuthorized), allowedChannel, chatChanFn)
}

// Options mocks base method.
func (m *MockEvebotCommand) Options() CommandOptions {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Options")
	ret0, _ := ret[0].(CommandOptions)
	return ret0
}

// Options indicates an expected call of Options.
func (mr *MockEvebotCommandMockRecorder) Options() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Options", reflect.TypeOf((*MockEvebotCommand)(nil).Options))
}
