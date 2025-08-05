package models

import "github.com/google/uuid"

type Repository interface {
	Open() error
	Close() error

	Put(data *Data) error
	GetByID(id uuid.UUID) (*Data, error)
	ListByPeriod(from, to int64) ([]Data, error)
}
