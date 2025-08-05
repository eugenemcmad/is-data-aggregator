package service

import (
	"errors"
	"github.com/google/uuid"
	"xis-data-aggregator/internal/models"
)

var (
	ErrNotFound = errors.New("not found")
	ErrCorrupt  = errors.New("corrupted data")
)

type DataService struct {
	repo models.Repository
}

func NewDataService(repo models.Repository) *DataService {
	return &DataService{repo: repo}
}

func (o *DataService) Put(data *models.Data) error {
	return o.repo.Put(data)
}

func (o *DataService) GetByID(id uuid.UUID) (*models.Data, error) {

	data, err := o.repo.GetByID(id)
	switch {
	case err != nil:
		return nil, err
	case data == nil:
		return nil, ErrNotFound
	case data.ID == uuid.Nil:
		return nil, ErrCorrupt
	}

	return data, nil
}

func (o *DataService) ListByPeriod(from, to int64) ([]models.Data, error) {
	data, err := o.repo.ListByPeriod(from, to)

	switch {
	case err != nil:
		return []models.Data{}, err
	case data == nil:
		return nil, ErrNotFound
	case len(data) == 0:
		return []models.Data{}, ErrNotFound
	}

	return data, nil
}
