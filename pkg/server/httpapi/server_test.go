package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/mymmrac/project-glynn/internal/mocks"
	"github.com/mymmrac/project-glynn/pkg/server"
	"github.com/mymmrac/project-glynn/pkg/uuid"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestServer_sendMassage(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mocks.NewMockRepository(ctrl)
	log, _ := test.NewNullLogger()
	service := server.NewService(m, log)
	srv := Server{
		service: service,
		log:     log,
	}
	roomID := uuid.New()
	vars := map[string]string{
		roomIDParameter: roomID.String(),
	}

	newMessage := server.ChatNewMessage{
		UserID: uuid.New(),
		Text:   "test",
	}
	messageBytes, err := json.Marshal(newMessage)
	assert.NoError(t, err)

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
