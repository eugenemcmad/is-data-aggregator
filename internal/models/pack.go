package models

import (
	"github.com/google/uuid"
)

// External data model for input
type Pack struct {
	ID        uuid.UUID // UUID RFC9562  (psql 16 bytes)
	Timestamp int64
	Data      []int
}
