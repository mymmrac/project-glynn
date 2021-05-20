// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mymmrac/project-glynn/pkg/repository (interfaces: Repository)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	message "github.com/mymmrac/project-glynn/pkg/data/message"
	user "github.com/mymmrac/project-glynn/pkg/data/user"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// GetMessageTime mocks base method.
func (m *MockRepository) GetMessageTime(arg0 uuid.UUID) (time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessageTime", arg0)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMessageTime indicates an expected call of GetMessageTime.
func (mr *MockRepositoryMockRecorder) GetMessageTime(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessageTime", reflect.TypeOf((*MockRepository)(nil).GetMessageTime), arg0)
}

// GetMessages mocks base method.
func (m *MockRepository) GetMessages(arg0 uuid.UUID, arg1 time.Time, arg2 uint) ([]message.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessages", arg0, arg1, arg2)
	ret0, _ := ret[0].([]message.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMessages indicates an expected call of GetMessages.
func (mr *MockRepositoryMockRecorder) GetMessages(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessages", reflect.TypeOf((*MockRepository)(nil).GetMessages), arg0, arg1, arg2)
}

// GetUsersFromIDs mocks base method.
func (m *MockRepository) GetUsersFromIDs(arg0 []uuid.UUID) ([]user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersFromIDs", arg0)
	ret0, _ := ret[0].([]user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsersFromIDs indicates an expected call of GetUsersFromIDs.
func (mr *MockRepositoryMockRecorder) GetUsersFromIDs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersFromIDs", reflect.TypeOf((*MockRepository)(nil).GetUsersFromIDs), arg0)
}

// IsRoomExist mocks base method.
func (m *MockRepository) IsRoomExist(arg0 uuid.UUID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsRoomExist", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsRoomExist indicates an expected call of IsRoomExist.
func (mr *MockRepositoryMockRecorder) IsRoomExist(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsRoomExist", reflect.TypeOf((*MockRepository)(nil).IsRoomExist), arg0)
}

// SaveMessage mocks base method.
func (m *MockRepository) SaveMessage(arg0 *message.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveMessage", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveMessage indicates an expected call of SaveMessage.
func (mr *MockRepositoryMockRecorder) SaveMessage(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveMessage", reflect.TypeOf((*MockRepository)(nil).SaveMessage), arg0)
}
