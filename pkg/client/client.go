package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mymmrac/project-glynn/pkg/data/chat"
	"github.com/mymmrac/project-glynn/pkg/server/httpapi"
	"github.com/mymmrac/project-glynn/pkg/uuid"
)

const (
	baseURL          = "%s/api/"
	messagesEndpoint = "rooms/%s/messages"
)

const updateInterval = 1 * time.Second

// Client manages connection to server
type Client struct {
	httpClient *http.Client
	host       string
	roomID     string
	running    chan struct{}
}

// NewClient creates new client with connection to specified host
func NewClient(host string) *Client {
	return &Client{
		httpClient: http.DefaultClient,
		host:       host,
	}
}

// StartChat begins to listen for new messages and reading to send message until an error occurs
func (c *Client) StartChat(roomID string) {
	c.roomID = roomID
	c.running = make(chan struct{}, 1)
	go c.readMessages()
	go c.sendMessages()
	<-c.running
}

func (c *Client) readMessages() {
	defer func() {
		c.running <- struct{}{}
	}()

	var lastMessageID *uuid.UUID
	url := fmt.Sprintf(baseURL+messagesEndpoint, c.host, c.roomID)

	for {
		reqURL := url
		if lastMessageID != nil {
			reqURL += fmt.Sprintf("?%s=%s", httpapi.LastMessageIDParameter, lastMessageID)
		}

		resp, err := c.httpClient.Get(reqURL)
		if err != nil {
			fmt.Printf("Oppss...\n%v", err)
			return
		}

		switch resp.StatusCode {
		case http.StatusOK:
		//	nothing
		case http.StatusNotFound:
			fmt.Printf("Room with id %q not found!", c.roomID)
			return
		default:
			fmt.Printf("Oppss... Status code: %d [%s]", resp.StatusCode, resp.Status)
			return
		}

		var cm chat.Messages
		err = json.NewDecoder(resp.Body).Decode(&cm)
		if err != nil {
			fmt.Printf("Oppss...\n%v", err)
			return
		}

		if err = resp.Body.Close(); err != nil {
			fmt.Printf("Oppss...\n%v", err)
			return
		}

		for _, m := range cm.Messages {
			fmt.Printf("%s [\033[33m%s\033[0m]: %s\n", m.Time.Local().Format(time.RFC822), cm.Usernames[m.UserID], m.Text)
		}

		l := len(cm.Messages)
		if l > 0 {
			lastMessageID = &cm.Messages[l-1].ID
		}

		time.Sleep(updateInterval)
	}
}

func (c *Client) sendMessages() {
	defer func() {
		c.running <- struct{}{}
	}()

	url := fmt.Sprintf(baseURL+messagesEndpoint, c.host, c.roomID)
	userID, err := uuid.Parse("506a43a4-25e2-4017-bc0c-90084d784958")
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Print("\u001B[F\u001B[2K")

		text = c.parseText(text)
		if text == "" {
			continue
		}

		newMessage := chat.NewMessage{
			UserID: userID,
			Text:   text,
		}

		byteSlice, err := json.Marshal(newMessage)
		if err != nil {
			fmt.Printf("Oppss...\n%v", err)
			return
		}
		body := bytes.NewReader(byteSlice)

		resp, err := c.httpClient.Post(url, "application/json; charset=UTF-8", body)
		if err != nil {
			fmt.Printf("Oppss...\n%v", err)
			return
		}
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Oppss...\n%v", err)
			return
		}

		if resp.StatusCode != http.StatusCreated {
			fmt.Printf("Oppss... Status code: %d [%s]", resp.StatusCode, resp.Status)
			return
		}
	}
}

func (c *Client) parseText(text string) string {
	text = strings.TrimSpace(text)
	text = strings.ToValidUTF8(text, "")
	return text
}
