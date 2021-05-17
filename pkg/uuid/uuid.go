package uuid

import "github.com/google/uuid"

const Regex = "\\b[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89aAbB][a-f0-9]{3}-[a-f0-9]{12}\\b"

type UUID = uuid.UUID

func New() UUID {
	return uuid.New()
}

func Parse(s string) (UUID, error) {
	return uuid.Parse(s)
}

func ToStrings(uuids []UUID) []string {
	uuidsStr := make([]string, len(uuids))
	for i, id := range uuids {
		uuidsStr[i] = id.String()
	}
	return uuidsStr
}
