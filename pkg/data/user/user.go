package user

import (
	"github.com/mymmrac/project-glynn/pkg/uuid"
)

// TODO docs
type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}
