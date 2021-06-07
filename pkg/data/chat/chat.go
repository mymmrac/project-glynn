package chat

import (
	"github.com/mymmrac/project-glynn/pkg/data/message"
	"github.com/mymmrac/project-glynn/pkg/uuid"
)

// Messages represents slice of messages with corresponding map of usernames
type Messages struct {
	Messages  []message.Message    `json:"messages"`  // Messages ordered by time
	Usernames map[uuid.UUID]string `json:"usernames"` // Usernames of users who sent messages
}

// NewMessage represents new message from users
type NewMessage struct {
	UserID uuid.UUID `json:"userID"` // UserID who sent message
	Text   string    `json:"text"`   // Text of sent message
}

// NewUser represents new user to be created
type NewUser struct {
	Username string `json:"username"` // Username of new user to be created
}
