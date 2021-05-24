package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/mymmrac/project-glynn/pkg/data/chat"
	"github.com/mymmrac/project-glynn/pkg/server/httpapi"
	"github.com/mymmrac/project-glynn/pkg/uuid"
)

const (
	baseURL          = "%s/api/"
	messagesEndpoint = "rooms/%s/messages"
	usersEndpoint    = "users"
)

const updateInterval = 1 * time.Second
const clearCurrentLine = "\u001B[F\u001B[2K"

var usernameRegex = regexp.MustCompile(`^[a-zA-Z]\w{2,}$`)

// Client manages connection to server
type Client struct {
	httpClient     *http.Client
	out            io.Writer
	in             io.Reader
	host           string
	roomID         string
	running        chan struct{}
	updateInterval time.Duration
}

// NewClient creates new client with connection to specified host
func NewClient(host string) *Client {
	return &Client{
		httpClient:     http.DefaultClient,
		host:           host,
		out:            os.Stdout,
		in:             os.Stdin,
		updateInterval: updateInterval,
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
			fmt.Fprintf(c.out, "Unable to make get request.\nError: %v\n", err)
			return
		}

		switch resp.StatusCode {
		case http.StatusOK:
		//	nothing
		case http.StatusNotFound:
			fmt.Fprintf(c.out, "Room with id %q not found.\n", c.roomID)
			return
		default:
			fmt.Fprintf(c.out, "Something went wrong.\nStatus code: %d [%s]\n", resp.StatusCode, resp.Status)
			return
		}

		var cm chat.Messages
		err = json.NewDecoder(resp.Body).Decode(&cm)
		if err != nil {
			fmt.Fprintf(c.out, "Unable to decode message.\nError: %v\n", err)
			return
		}

		if err = resp.Body.Close(); err != nil {
			fmt.Fprintf(c.out, "Unable to close response body.\nError: %v\n", err)
			return
		}

		// TODO move displaying to other func
		for _, m := range cm.Messages {
			fmt.Fprintf(c.out,
				"%s [\033[33m%s\033[0m]: %s\n", m.Time.Local().Format(time.RFC822), cm.Usernames[m.UserID], m.Text)
		}

		l := len(cm.Messages)
		if l > 0 {
			lastMessageID = &cm.Messages[l-1].ID
		}

		time.Sleep(c.updateInterval)
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

	scanner := bufio.NewScanner(c.in)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Fprint(c.out, clearCurrentLine)

		text = c.parseText(text)
		if text == "" {
			fmt.Fprintln(c.out, "Empty or non UTF-8 text can't be sent.")
			continue
		}

		newMessage := chat.NewMessage{
			UserID: userID,
			Text:   text,
		}

		byteSlice, err := json.Marshal(newMessage)
		if err != nil {
			fmt.Fprintf(c.out, "Unable to encode message.\nError: %v\n", err)
			return
		}
		body := bytes.NewReader(byteSlice)

		resp, err := c.httpClient.Post(url, "application/json; charset=UTF-8", body)
		if err != nil {
			fmt.Fprintf(c.out, "Unable to make post request.\nError: %v\n", err)
			return
		}
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(c.out, "Unable to close response body.\nError: %v\n", err)
			return
		}

		if resp.StatusCode != http.StatusCreated {
			fmt.Fprintf(c.out, "Something went wrong.\nStatus code: %d [%s]\n", resp.StatusCode, resp.Status)
			return
		}
	}
}

func (c *Client) parseText(text string) string {
	text = strings.ToValidUTF8(text, "")
	text = strings.TrimSpace(text)
	return text
}

func (c *Client) CreateUser(username string) {
	if !usernameRegex.MatchString(username) {
		fmt.Fprint(c.out, "Invalid username, must contain only [a-Z], [0-9] or '_', "+
			"starting from letter and at least 3 chars long.")
		return
	}

	newUser := chat.NewUser{
		Username: username,
	}

	byteSlice, err := json.Marshal(newUser)
	if err != nil {
		fmt.Fprintf(c.out, "Unable to encode user.\nError: %v\n", err)
		return
	}
	body := bytes.NewReader(byteSlice)

	url := fmt.Sprintf(baseURL+usersEndpoint, c.host)

	resp, err := c.httpClient.Post(url, "application/json; charset=UTF-8", body)
	if err != nil {
		fmt.Fprintf(c.out, "Unable to make post request.\nError: %v\n", err)
		return
	}
	if err := resp.Body.Close(); err != nil {
		fmt.Fprintf(c.out, "Unable to close response body.\nError: %v\n", err)
		return
	}

	if resp.StatusCode != http.StatusCreated {
		fmt.Fprintf(c.out, "Something went wrong.\nStatus code: %d [%s]\n", resp.StatusCode, resp.Status)
		return
	}

	var userIDBytes []byte
	_, err = resp.Body.Read(userIDBytes)
	if err != nil {
		fmt.Fprintf(c.out, "Unable to read response body.")
		return
	}

	err = ioutil.WriteFile("user.data", userIDBytes, 0600)
	if err != nil {
		fmt.Fprintf(c.out, "Unable to save user info.")
		return
	}

	fmt.Fprintf(c.out, "User created succsefully, now you can join rooms.")
}
