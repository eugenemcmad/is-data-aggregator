package models

import (
	"github.com/google/uuid"
	"xis-data-aggregator/pkg/utils"
)

// Data model for export
type Data struct {
	ID        uuid.UUID `json:"id"`
	Timestamp int64     `json:"ts"`
	Max       int       `json:"max"`
}

func MapPackToData(pack *Pack) (*Data, error) {
	data := Data{ID: pack.ID, Timestamp: pack.Timestamp}

	var err error
	data.Max, err = utils.GetMaxValue(pack.Data)

	return &data, err

}
