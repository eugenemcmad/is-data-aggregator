package models

import (
	"xis-data-aggregator/pkg/utils"

	"github.com/google/uuid"
)

// Data represents the processed and aggregated data model for export.
// This struct contains the essential information extracted from raw Pack data
// and is used for API responses and data storage.
type Data struct {
	ID        uuid.UUID `json:"id"`  // Unique identifier for the data record
	Timestamp int64     `json:"ts"`  // Unix timestamp when the data was recorded
	Max       int       `json:"max"` // Maximum value extracted from the original data array
}

// MapPackToData converts a Pack struct to a Data struct by extracting
// the maximum value from the Pack's data array and creating a simplified
// representation suitable for export and API responses.
//
// Parameters:
//   - pack: Pointer to the source Pack containing raw data
//
// Returns:
//   - *Data: Pointer to the converted Data struct
//   - error: Any error that occurred during the conversion process
func MapPackToData(pack *Pack) (*Data, error) {
	data := Data{ID: pack.ID, Timestamp: pack.Timestamp}

	var err error
	data.Max, err = utils.GetMaxValue(pack.Data)

	return &data, err

}
