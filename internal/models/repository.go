package models

import "github.com/google/uuid"

// Repository defines the interface for data persistence operations.
// This interface abstracts the storage layer and provides methods
// for storing, retrieving, and querying Data records.
type Repository interface {
	// Open initializes the repository connection and prepares it for use.
	// This method should be called before any other operations.
	Open() error

	// Close properly shuts down the repository connection and releases resources.
	// This method should be called when the repository is no longer needed.
	Close() error

	// Put stores a Data record in the repository.
	// If a record with the same ID already exists, it will be overwritten.
	//
	// Parameters:
	//   - data: Pointer to the Data struct to be stored
	//
	// Returns:
	//   - error: Any error that occurred during the storage operation
	Put(data *Data) error

	// GetByID retrieves a Data record by its unique identifier.
	//
	// Parameters:
	//   - id: UUID of the record to retrieve
	//
	// Returns:
	//   - *Data: Pointer to the retrieved Data struct, or nil if not found
	//   - error: Any error that occurred during the retrieval operation
	GetByID(id uuid.UUID) (*Data, error)

	// ListByPeriod retrieves all Data records within a specified time period.
	// The search is inclusive of both the 'from' and 'to' timestamps.
	//
	// Parameters:
	//   - from: Start timestamp (inclusive) for the search period
	//   - to: End timestamp (inclusive) for the search period
	//
	// Returns:
	//   - []Data: Slice of Data records found within the specified period
	//   - error: Any error that occurred during the search operation
	ListByPeriod(from, to int64) ([]Data, error)
}
