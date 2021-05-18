package user

import (
	"github.com/mymmrac/project-glynn/pkg/uuid"
)

// User represents info about chat participant
type User struct {
	ID       uuid.UUID `json:"id"`       // ID is a uniq identifier of user
	Username string    `json:"username"` // Username is name of user which can will be displayed
}
