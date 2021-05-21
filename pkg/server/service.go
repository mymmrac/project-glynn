package server

import (
	"errors"
	"fmt"
	"time"

	"github.com/mymmrac/project-glynn/pkg/data/chat"
	"github.com/mymmrac/project-glynn/pkg/data/message"
	"github.com/mymmrac/project-glynn/pkg/repository"
	"github.com/mymmrac/project-glynn/pkg/uuid"
	"github.com/sirupsen/logrus"
)

// MessageLimit limits amount of messages to be received
const MessageLimit uint = 20

var ErrorRoomNotFound = errors.New("no such room")

// Service manages all logic for api
type Service struct {
	messageRepo repository.MessageRepository
	userRepo    repository.UserRepository
	roomRepo    repository.RoomRepository
	log         *logrus.Logger
}

// NewService creates new Service with repository.Repository
func NewService(repo repository.Repository, log *logrus.Logger) *Service {
	return &Service{
		messageRepo: repo,
		userRepo:    repo,
		roomRepo:    repo,
		log:         log,
	}
}

// GetMessagesAfterTime returns chat.Messages after specified time
func (s *Service) GetMessagesAfterTime(roomID uuid.UUID, afterTime time.Time) (*chat.Messages, error) {
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

	return &chat.Messages{
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

// GetMessagesAfterMessage returns chat.Messages after specified message
func (s *Service) GetMessagesAfterMessage(roomID, lastMessageID uuid.UUID) (*chat.Messages, error) {
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

// GetMessagesLatest returns latest chat.Messages
func (s *Service) GetMessagesLatest(roomID uuid.UUID) (*chat.Messages, error) {
	cm, err := s.GetMessagesAfterTime(roomID, time.Time{})
	if err != nil {
		return nil, fmt.Errorf("latest messages: %w", err)
	}
	return cm, nil
}

// SendMessage saves message
func (s *Service) SendMessage(roomID uuid.UUID, newMessage chat.NewMessage) error {
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

// CheckRoom returns error if room not exist
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
