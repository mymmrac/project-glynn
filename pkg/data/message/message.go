package message

import (
	"fmt"
	"time"

	"github.com/mymmrac/project-glynn/pkg/uuid"
)

// TODO docs
type Message struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"userID"`
	RoomID uuid.UUID `json:"roomID"`
	Text   string    `json:"text"`
	Time   time.Time `json:"time"`
}

func (m *Message) String() string {
	return fmt.Sprintf("Message %s: '%s', at: %s , from: %s, room: %s", m.ID, m.Text, m.Time, m.UserID, m.RoomID)
}
