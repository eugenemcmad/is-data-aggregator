package models

import (
	"github.com/google/uuid"
)

// Pack represents the external data model for input processing.
// This struct contains raw data received from external sources
// and serves as the primary input format for the data aggregation system.
type Pack struct {
	ID        uuid.UUID // UUID RFC9562 (psql 16 bytes) - Unique identifier for the data pack
	Timestamp int64     // Unix timestamp indicating when the data was collected
	Data      []int     // Array of integer values representing the raw data points
}
