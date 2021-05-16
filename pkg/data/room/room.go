package room

import "github.com/google/uuid"

type Room struct {
	ID uuid.UUID `json:"id,omitempty"`
}
