package repository

import (
	"time"

	"github.com/mymmrac/project-glynn/pkg/data/message"
	"github.com/mymmrac/project-glynn/pkg/data/user"
	"github.com/mymmrac/project-glynn/pkg/uuid"
)

// TODO split to 3 repo
// TODO docs
type Repository interface {
	GetMessageTime(messageID uuid.UUID) (time.Time, error)
	GetMessages(roomID uuid.UUID, afterTime time.Time, limit uint) ([]message.Message, error)
	GetUsersFromIDs([]uuid.UUID) ([]user.User, error)
	SaveMessage(message *message.Message) error
	IsRoomExist(roomID uuid.UUID) (bool, error)
}
