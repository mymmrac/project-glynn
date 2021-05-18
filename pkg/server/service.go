package server

import (
	"errors"
	"fmt"
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
	messageRepo repository.MessageRepository
	userRepo    repository.UserRepository
	roomRepo    repository.RoomRepository
	log         *logrus.Logger
}

func NewService(repo repository.Repository, log *logrus.Logger) *Service {
	return &Service{
		messageRepo: repo,
		userRepo:    repo,
		roomRepo:    repo,
		log:         log,
	}
}

func (s *Service) GetMessagesAfterTime(roomID uuid.UUID, afterTime time.Time) (*ChatMessages, error) {
	if err := s.CheckRoom(roomID); err != nil {
		return nil, fmt.Errorf("messages after time: %w", err)
	}

	messages, err := s.messageRepo.GetMessages(roomID, afterTime, MessageLimit)
	if err != nil {
		return nil, fmt.Errorf("messages after time: %w", err)
	}

	if messages == nil {
		messages = make([]message.Message, 0)
	}

	ids := s.getUserIDsFromMessages(messages)
	users, err := s.userRepo.GetUsersFromIDs(ids)
	if err != nil {
		return nil, fmt.Errorf("messages after time: %w", err)
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

func (s Service) getUserIDsFromMessages(messages []message.Message) []uuid.UUID {
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

	return ids
}

func (s *Service) GetMessagesAfterMessage(roomID, lastMessageID uuid.UUID) (*ChatMessages, error) {
	msgTime, err := s.messageRepo.GetMessageTime(lastMessageID)
	if err != nil {
		return nil, fmt.Errorf("messages after message: %w", err)
	}

	cm, err := s.GetMessagesAfterTime(roomID, msgTime)
	if err != nil {
		return nil, fmt.Errorf("messages after message: %w", err)
	}
	return cm, nil
}

func (s *Service) GetMessagesLatest(roomID uuid.UUID) (*ChatMessages, error) {
	cm, err := s.GetMessagesAfterTime(roomID, time.Time{})
	if err != nil {
		return nil, fmt.Errorf("latest messages: %w", err)
	}
	return cm, nil
}

func (s *Service) SendMessage(roomID uuid.UUID, newMessage ChatNewMessage) error {
	if err := s.CheckRoom(roomID); err != nil {
		return fmt.Errorf("send message: %w", err)
	}
	// TODO user check
	// TODO text check

	msg := &message.Message{
		ID:     uuid.New(),
		UserID: newMessage.UserID,
		RoomID: roomID,
		Text:   newMessage.Text,
		Time:   time.Now(),
	}

	if err := s.messageRepo.SaveMessage(msg); err != nil {
		return fmt.Errorf("send message: %w", err)
	}
	return nil
}

func (s *Service) CheckRoom(roomID uuid.UUID) error {
	ok, err := s.roomRepo.IsRoomExist(roomID)
	if err != nil {
		return fmt.Errorf("check room: %w", err)
	}
	if !ok {
		return fmt.Errorf("check room: %w", ErrorRoomNotFound)
	}
	return nil
}
