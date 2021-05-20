package uuid

import "github.com/google/uuid"

// Regex for validating UUID
const Regex = "\\b[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[89aAbB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}\\b"

// UUID version 4 uuid
type UUID = uuid.UUID

// New generates new random UUID
func New() UUID {
	return uuid.New()
}

// Parse returns UUID from string if possible
func Parse(s string) (UUID, error) {
	return uuid.Parse(s)
}

// ToStrings returns slice of strings representations of UUID
func ToStrings(uuids []UUID) []string {
	uuidsStr := make([]string, len(uuids))
	for i, id := range uuids {
		uuidsStr[i] = id.String()
	}
	return uuidsStr
}
