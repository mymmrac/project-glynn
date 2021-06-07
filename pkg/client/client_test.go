package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/bmizerany/assert"
	"github.com/mymmrac/project-glynn/pkg/data/chat"
	"github.com/mymmrac/project-glynn/pkg/data/message"
	"github.com/mymmrac/project-glynn/pkg/uuid"
	"github.com/stretchr/testify/require"
)

func TestClient_parseText(t *testing.T) {
	c := Client{}

	type args struct {
		text string
	}
	type expected struct {
		text string
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name:     "empty",
			args:     args{text: ""},
			expected: expected{text: ""},
		},
		{
			name:     "ok",
			args:     args{text: "test text :)"},
			expected: expected{text: "test text :)"},
		},
		{
			name:     "spaces",
			args:     args{text: "     "},
			expected: expected{text: ""},
		},
		{
			name:     "spaces",
			args:     args{text: "     test  test    "},
			expected: expected{text: "test  test"},
		},
		{
			name:     "non uft-8",
			args:     args{text: "test \xf0\x90\xbc"},
			expected: expected{text: "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := c.parseText(tt.args.text)
			assert.Equal(t, tt.expected.text, actual)
		})
	}
}

func TestClient_readMessages(t *testing.T) {
	roomID := uuid.New()
	userID := uuid.New()
	messageTime := time.Unix(1621521072, 0).UTC()
	url := fmt.Sprintf(baseURL+messagesEndpoint, "", roomID)

	cm := chat.Messages{
		Messages: []message.Message{
			{ID: uuid.New(), UserID: userID, RoomID: roomID, Text: "test", Time: messageTime},
		},
		Usernames: map[uuid.UUID]string{
			userID: "test",
		},
	}

	runTimes := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, url, r.URL.Path)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		switch runTimes {
		case 0:
			err := json.NewEncoder(w).Encode(cm)
			require.NoError(t, err)
		case 1:
			err := json.NewEncoder(w).Encode(chat.Messages{
				Messages:  make([]message.Message, 0),
				Usernames: make(map[uuid.UUID]string),
			})
			require.NoError(t, err)
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		runTimes++
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	var outBuf bytes.Buffer
	c := &Client{
		httpClient:     server.Client(),
		host:           server.URL,
		roomID:         roomID.String(),
		running:        make(chan struct{}, 1),
		out:            &outBuf,
		updateInterval: 0,
	}

	c.readMessages()

	assert.Equal(t, 2, runTimes)

	expectedOutBuf := bytes.NewBufferString(
		"20 May 21 17:31 EEST [\u001B[33mtest\u001B[0m]: test\n" +
			"Something went wrong.\n" +
			"Status code: 400 [400 Bad Request]\n")
	assert.Equal(t, expectedOutBuf, &outBuf,
		fmt.Sprintf("expected: %q\n actual: %q", expectedOutBuf.String(), outBuf.String()))
}

func TestClient_sendMessages(t *testing.T) {
	roomID := uuid.New()
	url := fmt.Sprintf(baseURL+messagesEndpoint, "", roomID)

	runTimes := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, url, r.URL.Path)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		switch runTimes {
		case 0:
			w.WriteHeader(http.StatusCreated)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
		runTimes++
	}))

	inBuf := bytes.NewBufferString("test message\ntest message 2\n")

	var outBuf bytes.Buffer
	c := &Client{
		httpClient: server.Client(),
		host:       server.URL,
		roomID:     roomID.String(),
		running:    make(chan struct{}, 1),
		in:         inBuf,
		out:        &outBuf,
	}

	c.sendMessages()

	assert.Equal(t, 2, runTimes)

	expectedOutBuf := bytes.NewBufferString(
		clearCurrentLine + clearCurrentLine + "Something went wrong.\nStatus code: 400 [400 Bad Request]\n")
	assert.Equal(t, expectedOutBuf, &outBuf, fmt.Sprintf("%q", outBuf.String()))
}

func TestNewClient(t *testing.T) {
	host := "http://test.com"

	c := NewClient(host)

	assert.Equal(t, host, c.host)
	assert.Equal(t, os.Stdout, c.out)
	assert.Equal(t, os.Stdin, c.in)
	assert.Equal(t, http.DefaultClient, c.httpClient)
	assert.Equal(t, updateInterval, c.updateInterval)
}
