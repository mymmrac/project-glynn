package server

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/mymmrac/project-glynn/internal/mocks"
	"github.com/mymmrac/project-glynn/pkg/data/message"
	"github.com/mymmrac/project-glynn/pkg/data/user"
	"github.com/mymmrac/project-glynn/pkg/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var errAny = errors.New("any error")

func mockIsRoomExist(m *mocks.MockRepository, roomID gomock.Matcher, ok bool, err error) {
	m.EXPECT().
		IsRoomExist(roomID).
		Return(ok, err).
		Times(1)
}

func mockGetMessages(m *mocks.MockRepository,
	roomID gomock.Matcher, afterTime gomock.Matcher, messages []message.Message, err error) {
	m.EXPECT().
		GetMessages(roomID, afterTime, gomock.Eq(MessageLimit)).
		Return(messages, err).
		Times(1)
}

func mockGetUsersFromIDs(m *mocks.MockRepository, ids gomock.Matcher, users []user.User, err error) {
	m.EXPECT().
		GetUsersFromIDs(ids).
		Return(users, err).
		Times(1)
}

func getMessagesData(afterTime time.Time) (uuid.UUID, []user.User, []uuid.UUID, map[uuid.UUID]string, []message.Message) {
	roomID := uuid.New()

	users := make([]user.User, 3)
	ids := make([]uuid.UUID, len(users))
	usernames := make(map[uuid.UUID]string)
	for i := range users {
		users[i] = user.User{
			ID:       uuid.New(),
			Username: "user-" + strconv.Itoa(i),
		}
		ids[i] = users[i].ID
		usernames[users[i].ID] = users[i].Username
	}

	messages := make([]message.Message, 5)
	for i := range messages {
		messages[i] = message.Message{
			ID:     uuid.New(),
			UserID: users[i%len(users)].ID,
			RoomID: roomID,
			Text:   "message " + strconv.Itoa(i),
			Time:   afterTime.Add(time.Second * time.Duration(i+1)),
		}
	}

	return roomID, users, ids, usernames, messages
}

func TestService_GetMessagesAfterTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mocks.NewMockRepository(ctrl)
	service := NewService(m, logrus.StandardLogger())

	afterTime := time.Now()
	roomID, users, _, usernames, messages := getMessagesData(afterTime)

	t.Run("ok", func(t *testing.T) {
		mockIsRoomExist(m, gomock.Eq(roomID), true, nil)
		mockGetMessages(m, gomock.Eq(roomID), gomock.Eq(afterTime), messages, nil)
		mockGetUsersFromIDs(m, gomock.Any(), users, nil)

		actual, err := service.GetMessagesAfterTime(roomID, afterTime)
		assert.NoError(t, err)
		assert.Equal(t,
			&ChatMessages{
				Messages:  messages,
				Usernames: usernames,
			},
			actual)
	})

	t.Run("check room err", func(t *testing.T) {
		mockIsRoomExist(m, gomock.Eq(roomID), false, nil)

		actual, err := service.GetMessagesAfterTime(roomID, afterTime)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	t.Run("get messages err", func(t *testing.T) {
		mockIsRoomExist(m, gomock.Eq(roomID), true, nil)
		mockGetMessages(m, gomock.Eq(roomID), gomock.Eq(afterTime), nil, errAny)

		actual, err := service.GetMessagesAfterTime(roomID, afterTime)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	t.Run("get usernames err", func(t *testing.T) {
		mockIsRoomExist(m, gomock.Eq(roomID), true, nil)
		mockGetMessages(m, gomock.Eq(roomID), gomock.Eq(afterTime), messages, nil)
		mockGetUsersFromIDs(m, gomock.Any(), nil, errAny)

		actual, err := service.GetMessagesAfterTime(roomID, afterTime)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestService_getUserIDsFromMessages(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mocks.NewMockRepository(ctrl)
	service := NewService(m, logrus.StandardLogger())

	type args struct {
		messages []message.Message
	}
	ids := []uuid.UUID{
		uuid.New(),
		uuid.New(),
		uuid.New(),
	}
	tests := []struct {
		name     string
		args     args
		expected []uuid.UUID
	}{
		{
			name:     "nil",
			args:     struct{ messages []message.Message }{messages: nil},
			expected: []uuid.UUID{},
		},
		{
			name:     "empty",
			args:     struct{ messages []message.Message }{messages: []message.Message{}},
			expected: []uuid.UUID{},
		},
		{
			name: "one",
			args: struct{ messages []message.Message }{messages: []message.Message{
				{ID: uuid.UUID{}, UserID: ids[0], RoomID: uuid.UUID{}, Text: "", Time: time.Time{}},
			}},
			expected: []uuid.UUID{
				ids[0],
			},
		},
		{
			name: "many",
			args: struct{ messages []message.Message }{messages: []message.Message{
				{ID: uuid.UUID{}, UserID: ids[0], RoomID: uuid.UUID{}, Text: "", Time: time.Time{}},
				{ID: uuid.UUID{}, UserID: ids[1], RoomID: uuid.UUID{}, Text: "", Time: time.Time{}},
				{ID: uuid.UUID{}, UserID: ids[0], RoomID: uuid.UUID{}, Text: "", Time: time.Time{}},
				{ID: uuid.UUID{}, UserID: ids[2], RoomID: uuid.UUID{}, Text: "", Time: time.Time{}},
				{ID: uuid.UUID{}, UserID: ids[2], RoomID: uuid.UUID{}, Text: "", Time: time.Time{}},
			}},
			expected: []uuid.UUID{
				ids[0],
				ids[1],
				ids[2],
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := service.getUserIDsFromMessages(tt.args.messages)
			assert.ElementsMatch(t, tt.expected, actual)
		})
	}
}

func TestService_getUsernamesFromUserIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mocks.NewMockRepository(ctrl)
	service := NewService(m, logrus.StandardLogger())

	type args struct {
		ids []uuid.UUID
	}
	tests := []struct {
		name        string
		args        args
		users       []user.User
		expected    map[uuid.UUID]string
		expectedErr bool
	}{
		{
			name:        "nil",
			args:        struct{ ids []uuid.UUID }{ids: nil},
			users:       []user.User{},
			expected:    map[uuid.UUID]string{},
			expectedErr: false,
		},
		{
			name:        "empty",
			args:        struct{ ids []uuid.UUID }{ids: []uuid.UUID{}},
			users:       []user.User{},
			expected:    map[uuid.UUID]string{},
			expectedErr: false,
		},
		{
			name: "error",
			args: struct{ ids []uuid.UUID }{ids: []uuid.UUID{
				uuid.New(),
			}},
			users:       nil,
			expected:    map[uuid.UUID]string{},
			expectedErr: true,
		},
		{
			name: "ok",
			args: struct{ ids []uuid.UUID }{ids: []uuid.UUID{
				uuid.New(),
				uuid.New(),
			}},
			users: []user.User{
				{Username: "1"},
				{Username: "2"},
			},
			expected:    map[uuid.UUID]string{},
			expectedErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.expectedErr {
				err = errAny
			} else {
				for i, id := range tt.args.ids {
					tt.users[i].ID = id
					tt.expected[id] = tt.users[i].Username
				}
			}
			mockGetUsersFromIDs(m, gomock.Eq(tt.args.ids), tt.users, err)

			actual, err := service.getUsernamesFromUserIDs(tt.args.ids)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestService_CheckRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mocks.NewMockRepository(ctrl)
	service := NewService(m, logrus.StandardLogger())

	roomID := uuid.New()
	type expected struct {
		ok  bool
		err bool
	}
	tests := []struct {
		name     string
		expected expected
	}{
		{
			name: "ok",
			expected: expected{
				ok:  true,
				err: false,
			},
		},
		{
			name: "bad",
			expected: expected{
				ok:  false,
				err: false,
			},
		},
		{
			name: "err",
			expected: expected{
				ok:  false,
				err: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.expected.err {
				err = errAny
			}
			mockIsRoomExist(m, gomock.Eq(roomID), tt.expected.ok, err)

			err = service.CheckRoom(roomID)
			if tt.expected.ok {
				assert.NoError(t, err)
				return
			}
			assert.Error(t, err)
		})
	}
}

func TestService_SendMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mocks.NewMockRepository(ctrl)
	service := NewService(m, logrus.StandardLogger())

	type args struct {
		roomID     uuid.UUID
		newMessage ChatNewMessage
	}
	type expected struct {
		roomExistErr bool
		err          bool
	}

	roomID := uuid.New()

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "ok",
			args: args{
				roomID: roomID,
				newMessage: ChatNewMessage{
					UserID: uuid.New(),
					Text:   "Test",
				},
			},
			expected: expected{
				roomExistErr: false,
				err:          false,
			},
		},
		{
			name: "err room",
			args: args{
				roomID: roomID,
				newMessage: ChatNewMessage{
					UserID: uuid.New(),
					Text:   "Test",
				},
			},
			expected: expected{
				roomExistErr: true,
				err:          true,
			},
		},
		{
			name: "err",
			args: args{
				roomID: roomID,
				newMessage: ChatNewMessage{
					UserID: uuid.New(),
					Text:   "Test",
				},
			},
			expected: expected{
				roomExistErr: false,
				err:          true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			count := 1
			if tt.expected.roomExistErr {
				err = errAny
				count = 0
			}
			mockIsRoomExist(m, gomock.Eq(roomID), !tt.expected.roomExistErr, err)

			if tt.expected.err {
				err = errAny
			}
			m.EXPECT().
				SaveMessage(gomock.Any()).
				Return(err).
				Times(count)

			err = service.SendMessage(tt.args.roomID, tt.args.newMessage)
			if tt.expected.err {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestService_GetMessagesLatest(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mocks.NewMockRepository(ctrl)
	service := NewService(m, logrus.StandardLogger())

	afterTime := time.Now()
	roomID, users, _, usernames, messages := getMessagesData(afterTime)

	type expected struct {
		chatMessages *ChatMessages
		err          bool
	}
	tests := []struct {
		name     string
		expected expected
	}{
		{
			name: "ok",
			expected: expected{
				chatMessages: &ChatMessages{
					Messages:  messages,
					Usernames: usernames,
				},
				err: false,
			},
		},
		{
			name: "err",
			expected: expected{
				chatMessages: nil,
				err:          true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockIsRoomExist(m, gomock.Eq(roomID), true, nil)
			var err error
			if tt.expected.err {
				err = errAny
			}
			mockGetMessages(m, gomock.Eq(roomID), gomock.Any(), messages, err)
			if !tt.expected.err {
				mockGetUsersFromIDs(m, gomock.Any(), users, nil)
			}

			actual, err := service.GetMessagesLatest(roomID)
			if tt.expected.err {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.chatMessages, actual)
		})
	}
}

func TestService_GetMessagesAfterMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mocks.NewMockRepository(ctrl)
	service := NewService(m, logrus.StandardLogger())

	afterTime := time.Now()
	afterMessageID := uuid.New()
	roomID, users, _, usernames, messages := getMessagesData(afterTime)

	type expected struct {
		chatMessages   *ChatMessages
		messageTimeErr bool
		err            bool
	}
	tests := []struct {
		name     string
		expected expected
	}{
		{
			name: "ok",
			expected: expected{
				chatMessages: &ChatMessages{
					Messages:  messages,
					Usernames: usernames,
				},
				messageTimeErr: false,
				err:            false,
			},
		},
		{
			name: "err",
			expected: expected{
				chatMessages:   nil,
				messageTimeErr: false,
				err:            true,
			},
		},
		{
			name: "err message time",
			expected: expected{
				chatMessages:   nil,
				messageTimeErr: true,
				err:            true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.expected.messageTimeErr {
				err = errAny
			}
			m.EXPECT().
				GetMessageTime(gomock.Eq(afterMessageID)).
				Return(afterTime, err).
				Times(1)

			if !tt.expected.messageTimeErr {
				mockIsRoomExist(m, gomock.Eq(roomID), true, nil)

				err = nil
				if tt.expected.err {
					err = errAny
				}
				mockGetMessages(m, gomock.Eq(roomID), gomock.Any(), messages, err)
				if !tt.expected.err {
					mockGetUsersFromIDs(m, gomock.Any(), users, nil)
				}
			}

			actual, err := service.GetMessagesAfterMessage(roomID, afterMessageID)
			if tt.expected.err || tt.expected.messageTimeErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.chatMessages, actual)
		})
	}
}
