package uuid

import "github.com/google/uuid"

type UUID = uuid.UUID

// TODO simplify
const Regex = "\\b[0-9a-f]{8}\\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\\b[0-9a-f]{12}\\b"

func New() UUID {
	return uuid.New()
}

func Parse(s string) (UUID, error) {
	return uuid.Parse(s)
}
