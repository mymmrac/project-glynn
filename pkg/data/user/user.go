package user

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id,omitempty"`
	Username string    `json:"username"`
}
