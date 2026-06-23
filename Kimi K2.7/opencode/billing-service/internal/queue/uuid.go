package queue

import "github.com/google/uuid"

// UUID alias for internal use.
type UUID = uuid.UUID

func parseUUIDInternal(s string) (UUID, error) {
	return uuid.Parse(s)
}
