// Code generated by MockGen. DO NOT EDIT.
// Source: provider.go

// Package mock_chatservice is a generated GoMock package.
package chatservice

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	chatmodels "gitlab.unanet.io/devops/eve-bot/internal/chatservice/chatmodels"
)

// MockProvider is a mock of Provider interface
type MockProvider struct {
	ctrl     *gomock.Controller
	recorder *MockProviderMockRecorder
}

// MockProviderMockRecorder is the mock recorder for MockProvider
type MockProviderMockRecorder struct {
	mock *MockProvider
}

// NewMockProvider creates a new mock instance
func NewMockProvider(ctrl *gomock.Controller) *MockProvider {
	mock := &MockProvider{ctrl: ctrl}
	mock.recorder = &MockProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockProvider) EXPECT() *MockProviderMockRecorder {
	return m.recorder
}

// GetChannelInfo mocks base method
func (m *MockProvider) GetChannelInfo(ctx context.Context, channelID string) (chatmodels.Channel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChannelInfo", ctx, channelID)
	ret0, _ := ret[0].(chatmodels.Channel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChannelInfo indicates an expected call of GetChannelInfo
func (mr *MockProviderMockRecorder) GetChannelInfo(ctx, channelID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChannelInfo", reflect.TypeOf((*MockProvider)(nil).GetChannelInfo), ctx, channelID)
}

// PostMessage mocks base method
func (m *MockProvider) PostMessage(ctx context.Context, msg, channel string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostMessage", ctx, msg, channel)
	ret0, _ := ret[0].(string)
	return ret0
}

// PostMessage indicates an expected call of PostMessage
func (mr *MockProviderMockRecorder) PostMessage(ctx, msg, channel interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostMessage", reflect.TypeOf((*MockProvider)(nil).PostMessage), ctx, msg, channel)
}

// PostMessageThread mocks base method
func (m *MockProvider) PostMessageThread(ctx context.Context, msg, channel, ts string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostMessageThread", ctx, msg, channel, ts)
	ret0, _ := ret[0].(string)
	return ret0
}

// PostMessageThread indicates an expected call of PostMessageThread
func (mr *MockProviderMockRecorder) PostMessageThread(ctx, msg, channel, ts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostMessageThread", reflect.TypeOf((*MockProvider)(nil).PostMessageThread), ctx, msg, channel, ts)
}

// ErrorNotification mocks base method
func (m *MockProvider) ErrorNotification(ctx context.Context, user, channel string, err error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ErrorNotification", ctx, user, channel, err)
}

// ErrorNotification indicates an expected call of ErrorNotification
func (mr *MockProviderMockRecorder) ErrorNotification(ctx, user, channel, err interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ErrorNotification", reflect.TypeOf((*MockProvider)(nil).ErrorNotification), ctx, user, channel, err)
}

// ErrorNotificationThread mocks base method
func (m *MockProvider) ErrorNotificationThread(ctx context.Context, user, channel, ts string, err error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ErrorNotificationThread", ctx, user, channel, ts, err)
}

// ErrorNotificationThread indicates an expected call of ErrorNotificationThread
func (mr *MockProviderMockRecorder) ErrorNotificationThread(ctx, user, channel, ts, err interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ErrorNotificationThread", reflect.TypeOf((*MockProvider)(nil).ErrorNotificationThread), ctx, user, channel, ts, err)
}

// UserNotificationThread mocks base method
func (m *MockProvider) UserNotificationThread(ctx context.Context, msg, user, channel, ts string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UserNotificationThread", ctx, msg, user, channel, ts)
}

// UserNotificationThread indicates an expected call of UserNotificationThread
func (mr *MockProviderMockRecorder) UserNotificationThread(ctx, msg, user, channel, ts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserNotificationThread", reflect.TypeOf((*MockProvider)(nil).UserNotificationThread), ctx, msg, user, channel, ts)
}

// DeploymentNotificationThread mocks base method
func (m *MockProvider) DeploymentNotificationThread(ctx context.Context, msg, user, channel, ts string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeploymentNotificationThread", ctx, msg, user, channel, ts)
}

// DeploymentNotificationThread indicates an expected call of DeploymentNotificationThread
func (mr *MockProviderMockRecorder) DeploymentNotificationThread(ctx, msg, user, channel, ts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeploymentNotificationThread", reflect.TypeOf((*MockProvider)(nil).DeploymentNotificationThread), ctx, msg, user, channel, ts)
}

// GetUser mocks base method
func (m *MockProvider) GetUser(ctx context.Context, user string) (*chatmodels.ChatUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", ctx, user)
	ret0, _ := ret[0].(*chatmodels.ChatUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser
func (mr *MockProviderMockRecorder) GetUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockProvider)(nil).GetUser), ctx, user)
}

// PostLinkMessageThread mocks base method
func (m *MockProvider) PostLinkMessageThread(ctx context.Context, msg, user, channel, ts string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PostLinkMessageThread", ctx, msg, user, channel, ts)
}

// PostLinkMessageThread indicates an expected call of PostLinkMessageThread
func (mr *MockProviderMockRecorder) PostLinkMessageThread(ctx, msg, user, channel, ts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostLinkMessageThread", reflect.TypeOf((*MockProvider)(nil).PostLinkMessageThread), ctx, msg, user, channel, ts)
}

// ShowResultsMessageThread mocks base method
func (m *MockProvider) ShowResultsMessageThread(ctx context.Context, msg, user, channel, ts string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowResultsMessageThread", ctx, msg, user, channel, ts)
}

// ShowResultsMessageThread indicates an expected call of ShowResultsMessageThread
func (mr *MockProviderMockRecorder) ShowResultsMessageThread(ctx, msg, user, channel, ts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowResultsMessageThread", reflect.TypeOf((*MockProvider)(nil).ShowResultsMessageThread), ctx, msg, user, channel, ts)
}

// ReleaseResultsMessageThread mocks base method
func (m *MockProvider) ReleaseResultsMessageThread(ctx context.Context, msg, user, channel, ts string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ReleaseResultsMessageThread", ctx, msg, user, channel, ts)
}

// ReleaseResultsMessageThread indicates an expected call of ReleaseResultsMessageThread
func (mr *MockProviderMockRecorder) ReleaseResultsMessageThread(ctx, msg, user, channel, ts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReleaseResultsMessageThread", reflect.TypeOf((*MockProvider)(nil).ReleaseResultsMessageThread), ctx, msg, user, channel, ts)
}
