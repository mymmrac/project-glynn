package server

import (
	"errors"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/mymmrac/project-glynn/pkg/data/message"
	"github.com/mymmrac/project-glynn/pkg/repository"
	"github.com/sirupsen/logrus"
)

const MessageLimit uint = 20

var ErrorRoomNotFound = errors.New("no such room")

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

func (s *Service) GetMessagesAfterTime(roomID uuid.UUID, afterTime time.Time) ([]message.Message, map[uuid.UUID]string, error) {
	if err := s.CheckRoom(roomID); err != nil {
		return nil, nil, err
	}
	// TODO user check

	messages, err := s.repository.GetMessages(roomID, afterTime, MessageLimit)
	if err != nil {
		return nil, nil, err
	}

	if messages == nil {
		messages = make([]message.Message, 0)
	}

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Time.Before(messages[j].Time)
	})

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
		return nil, nil, err
	}

	usernames := make(map[uuid.UUID]string)
	for _, user := range users {
		usernames[user.ID] = user.Username
	}

	return messages, usernames, nil
}

func (s *Service) GetMessagesAfterMessage(roomID, lastMessageID uuid.UUID) ([]message.Message, map[uuid.UUID]string, error) {
	msgTime, err := s.repository.GetMessageTime(lastMessageID)
	if err != nil {
		return nil, nil, err
	}

	return s.GetMessagesAfterTime(roomID, msgTime)
}

func (s *Service) GetMessagesLatest(roomID uuid.UUID) ([]message.Message, map[uuid.UUID]string, error) {
	return s.GetMessagesAfterTime(roomID, time.Time{})
}

func (s *Service) SendMessage(roomID, userID uuid.UUID, text string) error {
	if err := s.CheckRoom(roomID); err != nil {
		return err
	}
	// TODO user check

	msg := &message.Message{
		ID:     uuid.New(),
		UserID: userID,
		RoomID: roomID,
		Text:   text,
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
