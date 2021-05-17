package server

import (
	"errors"
	"time"

	"github.com/mymmrac/project-glynn/pkg/data/message"
	"github.com/mymmrac/project-glynn/pkg/repository"
	"github.com/mymmrac/project-glynn/pkg/uuid"
	"github.com/sirupsen/logrus"
)

const MessageLimit uint = 20

var ErrorRoomNotFound = errors.New("no such room")

type ChatMessages struct {
	Messages  []message.Message    `json:"messages"`
	Usernames map[uuid.UUID]string `json:"usernames"`
}

type ChatNewMessage struct {
	UserID uuid.UUID `json:"userID"`
	Text   string    `json:"text"`
}

type Service struct {
	repository repository.Repository
	log        *logrus.Logger
}

func NewService(repo repository.Repository, log *logrus.Logger) *Service {
	return &Service{
		repository: repo,
		log:        log,
	}
}

func (s *Service) GetMessagesAfterTime(roomID uuid.UUID, afterTime time.Time) (*ChatMessages, error) {
	if err := s.CheckRoom(roomID); err != nil {
		return nil, err
	}

	messages, err := s.repository.GetMessages(roomID, afterTime, MessageLimit)
	if err != nil {
		return nil, err
	}

	if messages == nil {
		messages = make([]message.Message, 0)
	}

	// TODO move to own func
	idsMap := make(map[uuid.UUID]struct{})
	for _, msg := range messages {
		idsMap[msg.UserID] = struct{}{}
	}

	ids := make([]uuid.UUID, len(idsMap))
	i := 0
	for id := range idsMap {
		ids[i] = id
		i++
	}

	users, err := s.repository.GetUsersFromIDs(ids)
	if err != nil {
		return nil, err
	}

	usernames := make(map[uuid.UUID]string)
	for _, user := range users {
		usernames[user.ID] = user.Username
	}

	return &ChatMessages{
		Messages:  messages,
		Usernames: usernames,
	}, nil
}

func (s *Service) GetMessagesAfterMessage(roomID, lastMessageID uuid.UUID) (*ChatMessages, error) {
	msgTime, err := s.repository.GetMessageTime(lastMessageID)
	if err != nil {
		return nil, err
	}

	return s.GetMessagesAfterTime(roomID, msgTime)
}

func (s *Service) GetMessagesLatest(roomID uuid.UUID) (*ChatMessages, error) {
	return s.GetMessagesAfterTime(roomID, time.Time{})
}

func (s *Service) SendMessage(roomID uuid.UUID, newMessage ChatNewMessage) error {
	if err := s.CheckRoom(roomID); err != nil {
		return err
	}
	// TODO user check
	// TODO test check

	msg := &message.Message{
		ID:     uuid.New(),
		UserID: newMessage.UserID,
		RoomID: roomID,
		Text:   newMessage.Text,
		Time:   time.Now(),
	}

	return s.repository.SaveMessage(msg)
}

func (s *Service) CheckRoom(roomID uuid.UUID) error {
	ok, err := s.repository.IsRoomExist(roomID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrorRoomNotFound
	}
	return nil
}
