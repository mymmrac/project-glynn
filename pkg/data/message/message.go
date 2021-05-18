package message

import (
	"time"

	"github.com/mymmrac/project-glynn/pkg/uuid"
)

// Message represents one message from user in one room
type Message struct {
	ID     uuid.UUID `json:"id"`     // ID is a uniq identifier of message
	UserID uuid.UUID `json:"userID"` // UserID is an id of user who sent message
	RoomID uuid.UUID `json:"roomID"` // RoomID is an id of room in which message was sent
	Text   string    `json:"text"`   // Text is actual text that user sent
	Time   time.Time `json:"time"`   // Time when massage was sent
}
