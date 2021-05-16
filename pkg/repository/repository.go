package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/mymmrac/project-glynn/pkg/data/message"
	"github.com/mymmrac/project-glynn/pkg/data/user"
)

type Repository interface {
	GetMessageTime(messageID uuid.UUID) (time.Time, error)
	GetMessages(roomID uuid.UUID, afterTime time.Time, limit uint) ([]message.Message, error)
	GetUsersFromIDs([]uuid.UUID) ([]user.User, error)
	SaveMessage(message *message.Message) error
	IsRoomExist(roomID uuid.UUID) (bool, error)
}
