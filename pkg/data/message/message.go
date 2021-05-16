package message

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID     uuid.UUID `json:"id,omitempty"`
	UserID uuid.UUID `json:"userID"`
	RoomID uuid.UUID `json:"roomID"`
	Text   string    `json:"text"`
	Time   time.Time `json:"time,omitempty"`
}

func (m *Message) String() string {
	return fmt.Sprintf("Message %s: '%s', at: %s , from: %s, room: %s", m.ID, m.Text, m.Time, m.UserID, m.RoomID)
}
