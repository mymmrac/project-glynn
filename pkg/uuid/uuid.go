package uuid

import "github.com/google/uuid"

// Regex for matching uuid version 4
const Regex = "\\b[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89aAbB][a-f0-9]{3}-[a-f0-9]{12}\\b"

// UUID is a 16 byte Universal Unique Identifier as defined in RFC 4122
type UUID = uuid.UUID

// New returns new random UUID
func New() UUID {
	return uuid.New()
}

// Parse converts string into UUID
func Parse(s string) (UUID, error) {
	return uuid.Parse(s)
}

// ToStrings converts slice of UUID's to slice of strings
func ToStrings(uuids []UUID) []string {
	uuidsStr := make([]string, len(uuids))
	for i, id := range uuids {
		uuidsStr[i] = id.String()
	}
	return uuidsStr
}
