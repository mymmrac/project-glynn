package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/mymmrac/project-glynn/internal/mocks"
	"github.com/mymmrac/project-glynn/pkg/data/message"
	"github.com/mymmrac/project-glynn/pkg/data/user"
	"github.com/mymmrac/project-glynn/pkg/server"
	"github.com/mymmrac/project-glynn/pkg/uuid"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	m       *mocks.MockRepository
	service *server.Service
	srv     Server
	roomID  uuid.UUID
	vars    map[string]string
)

func setup(t *testing.T) {
	ctrl := gomock.NewController(t)
	m = mocks.NewMockRepository(ctrl)
	log, _ := test.NewNullLogger()
	service = server.NewService(m, log)
	srv = Server{
		service: service,
		log:     log,
	}
	roomID = uuid.New()
	vars = map[string]string{
		roomIDParameter: roomID.String(),
	}
}

func getTestData() (messages []message.Message, users []user.User, usernames map[uuid.UUID]string) {
	messages = []message.Message{
		{
			ID:     uuid.New(),
			RoomID: roomID,
			Text:   "message 1",
			Time:   time.Now().Round(0),
		},
		{
			ID:     uuid.New(),
			RoomID: roomID,
			Text:   "message 2",
			Time:   time.Now().Round(0),
		},
	}

	users = []user.User{
		{
			ID:       uuid.New(),
			Username: "Test",
		},
		{
			ID:       uuid.New(),
			Username: "TestOK",
		},
	}

	usernames = make(map[uuid.UUID]string)
	for i := range messages {
		j := i % len(users)
		messages[i].UserID = users[j].ID
		usernames[users[j].ID] = users[j].Username
	}

	return
}

func TestServer_sendMassage(t *testing.T) {
	setup(t)

	newMessage := server.ChatNewMessage{
		UserID: uuid.New(),
		Text:   "test",
	}
	messageBytes, err := json.Marshal(newMessage)
	require.NoError(t, err)

	reqNilBody := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/rooms/%s/messages", roomID), nil)

	t.Run("ok", func(t *testing.T) {
		mocks.MockIsRoomExist(m, gomock.Eq(roomID), true, nil)
		mocks.MockSaveMessage(m, gomock.Any(), nil, 1)

		req := httptest.NewRequest(http.MethodPost,
			fmt.Sprintf("/api/rooms/%s/messages", roomID),
			bytes.NewReader(messageBytes))
		vars := map[string]string{
			roomIDParameter: roomID.String(),
		}
		req = mux.SetURLVars(req, vars)

		rr := httptest.NewRecorder()
		srv.sendMassage()(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	t.Run("no room id", func(t *testing.T) {
		rr := httptest.NewRecorder()
		srv.sendMassage()(rr, reqNilBody)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("bad room id", func(t *testing.T) {
		req := reqNilBody
		vars := map[string]string{
			roomIDParameter: "test",
		}
		req = mux.SetURLVars(req, vars)

		rr := httptest.NewRecorder()
		srv.sendMassage()(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("decode message", func(t *testing.T) {
		req := reqNilBody
		req = mux.SetURLVars(req, vars)

		rr := httptest.NewRecorder()
		srv.sendMassage()(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("save err", func(t *testing.T) {
		mocks.MockIsRoomExist(m, gomock.Eq(roomID), true, nil)
		mocks.MockSaveMessage(m, gomock.Any(), errors.New("error"), 1)

		req := httptest.NewRequest(http.MethodPost,
			fmt.Sprintf("/api/rooms/%s/messages", roomID),
			bytes.NewReader(messageBytes))

		req = mux.SetURLVars(req, vars)

		rr := httptest.NewRecorder()
		srv.sendMassage()(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestServer_getMessages(t *testing.T) {
	setup(t)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/rooms/%s/messages", roomID), nil)
	req = mux.SetURLVars(req, vars)

	messages, users, usernames := getTestData()
	expected := &server.ChatMessages{
		Messages:  messages,
		Usernames: usernames,
	}
	var actual *server.ChatMessages

	withBadRequest := func(r *http.Request) {
		rr := httptest.NewRecorder()
		srv.getMessages()(rr, r)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	}

	t.Run("ok", func(t *testing.T) {
		mocks.MockIsRoomExist(m, gomock.Eq(roomID), true, nil)
		mocks.MockGetMessages(m, gomock.Eq(roomID), gomock.Any(), gomock.Eq(server.MessageLimit), messages, nil)
		mocks.MockGetUsersFromIDs(m, gomock.Any(), users, nil)

		rr := httptest.NewRecorder()
		srv.getMessages()(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		err := json.NewDecoder(rr.Body).Decode(&actual)
		assert.NoError(t, err)

		assert.Equal(t, expected, actual)
	})

	t.Run("ok last message", func(t *testing.T) {
		lastMessageID := uuid.New()
		afterTime := time.Now().Round(0)

		mocks.MockGetMessageTime(m, gomock.Eq(lastMessageID), afterTime, nil)
		mocks.MockIsRoomExist(m, gomock.Eq(roomID), true, nil)
		mocks.MockGetMessages(m, gomock.Eq(roomID), gomock.Eq(afterTime), gomock.Eq(server.MessageLimit), messages, nil)
		mocks.MockGetUsersFromIDs(m, gomock.Any(), users, nil)

		reqLastMessage := httptest.NewRequest(http.MethodGet,
			fmt.Sprintf("/api/rooms/%s/messages?%s=%s", roomID, lastMessageIDParameter, lastMessageID),
			nil)
		reqLastMessage = mux.SetURLVars(reqLastMessage, vars)

		rr := httptest.NewRecorder()
		srv.getMessages()(rr, reqLastMessage)

		assert.Equal(t, http.StatusOK, rr.Code)

		err := json.NewDecoder(rr.Body).Decode(&actual)
		assert.NoError(t, err)

		assert.Equal(t, expected, actual)
	})

	t.Run("bad last message id", func(t *testing.T) {
		reqLastMessage := httptest.NewRequest(http.MethodGet,
			fmt.Sprintf("/api/rooms/%s/messages?%s=%s", roomID, lastMessageIDParameter, "bad_id"),
			nil)
		reqLastMessage = mux.SetURLVars(reqLastMessage, vars)

		withBadRequest(reqLastMessage)
	})

	t.Run("no room id", func(t *testing.T) {
		reqNoRoomID := mux.SetURLVars(req, nil)

		withBadRequest(reqNoRoomID)
	})

	t.Run("bad room id", func(t *testing.T) {
		varsBad := map[string]string{
			roomIDParameter: "test",
		}
		reqBadRoomID := mux.SetURLVars(req, varsBad)

		withBadRequest(reqBadRoomID)
	})

	t.Run("get messages err", func(t *testing.T) {
		mocks.MockIsRoomExist(m, gomock.Eq(roomID), true, nil)
		mocks.MockGetMessages(m, gomock.Eq(roomID), gomock.Any(), gomock.Eq(server.MessageLimit), messages, errors.New(""))

		withBadRequest(req)
	})

	t.Run("room not found err", func(t *testing.T) {
		mocks.MockIsRoomExist(m, gomock.Eq(roomID), false, nil)

		rr := httptest.NewRecorder()
		srv.getMessages()(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}

func TestServer_routes(t *testing.T) {
	setup(t)

	srv.routes()

	type args struct {
		method string
		url    string
	}
	type expected struct {
		handler http.Handler
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "get messages",
			args: args{
				method: http.MethodGet,
				url:    fmt.Sprintf("/api/rooms/%s/messages", roomID),
			},
			expected: expected{
				handler: srv.getMessages(),
			},
		},
		{
			name: "get messages with last id",
			args: args{
				method: http.MethodGet,
				url:    fmt.Sprintf("/api/rooms/%s/messages?%s=%s", roomID, lastMessageIDParameter, uuid.New()),
			},
			expected: expected{
				handler: srv.getMessages(),
			},
		},
		{
			name: "send messages",
			args: args{
				method: http.MethodPost,
				url:    fmt.Sprintf("/api/rooms/%s/messages", roomID),
			},
			expected: expected{
				handler: srv.sendMassage(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.args.method, tt.args.url, nil)

			rm := &mux.RouteMatch{}
			srv.router.Match(request, rm)

			v1 := reflect.ValueOf(rm.Handler)
			v2 := reflect.ValueOf(tt.expected.handler)
			assert.Equal(t, v2.Pointer(), v1.Pointer(), "unexpected router")
		})
	}
}

func TestNewServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mocks.NewMockRepository(ctrl)
	log, _ := test.NewNullLogger()
	service := server.NewService(m, log)

	srv := NewServer(service, log)

	assert.Equal(t, log, srv.log)
	assert.Equal(t, service, srv.service)
}
