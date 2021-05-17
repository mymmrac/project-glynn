package room

import "github.com/mymmrac/project-glynn/pkg/uuid"

// Room represents chat room
type Room struct {
	ID uuid.UUID `json:"id"` // ID is a uniq identifier of room
}
