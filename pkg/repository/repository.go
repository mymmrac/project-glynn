package repository

import (
	"time"

	"github.com/mymmrac/project-glynn/pkg/data/message"
	"github.com/mymmrac/project-glynn/pkg/data/user"
	"github.com/mymmrac/project-glynn/pkg/uuid"
)

// Repository manages data related to messages, users and rooms
type Repository interface {
	MessageRepository
	UserRepository
	RoomRepository
}

// MessageRepository manages data related to messages
type MessageRepository interface {
	// GetMessageTime returns time when massage was sent by its id
	GetMessageTime(messageID uuid.UUID) (time.Time, error)

	// GetMessages returns limited amount of messages from specified room and after specified time
	GetMessages(roomID uuid.UUID, afterTime time.Time, limit uint) ([]message.Message, error)

	// SaveMessage saves given massage
	SaveMessage(message *message.Message) error
}

// UserRepository manages data related to users
type UserRepository interface {
	// GetUsersFromIDs returns slice of users by their ids
	GetUsersFromIDs([]uuid.UUID) ([]user.User, error)
}

// RoomRepository manages data related to rooms
type RoomRepository interface {
	// IsRoomExist checks if room exist
	IsRoomExist(roomID uuid.UUID) (bool, error)
}
