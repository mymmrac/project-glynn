package mocks

import (
	"time"

	"github.com/golang/mock/gomock"
	"github.com/mymmrac/project-glynn/pkg/data/message"
	"github.com/mymmrac/project-glynn/pkg/data/user"
)

func MockIsRoomExist(m *MockRepository, roomID gomock.Matcher, ok bool, err error) {
	m.EXPECT().
		IsRoomExist(roomID).
		Return(ok, err).
		Times(1)
}

func MockGetMessages(m *MockRepository,
	roomID, afterTime,
	messageLimit gomock.Matcher, messages []message.Message, err error) {
	m.EXPECT().
		GetMessages(roomID, afterTime, messageLimit).
		Return(messages, err).
		Times(1)
}

func MockGetUsersFromIDs(m *MockRepository, ids gomock.Matcher, users []user.User, err error) {
	m.EXPECT().
		GetUsersFromIDs(ids).
		Return(users, err).
		Times(1)
}

func MockGetMessageTime(m *MockRepository, afterMessageID gomock.Matcher, afterTime time.Time, err error) {
	m.EXPECT().
		GetMessageTime(afterMessageID).
		Return(afterTime, err).
		Times(1)
}

func MockSaveMessage(m *MockRepository, msg gomock.Matcher, err error, times int) {
	m.EXPECT().
		SaveMessage(msg).
		Return(err).
		Times(times)
}
