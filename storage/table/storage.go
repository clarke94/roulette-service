package table

import (
	"context"

	"github.com/clarke94/roulette-service/internal/pkg/table"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Storage provides a Storage layer.
type Storage struct {
	DB *gorm.DB
}

// New initializes Storage.
func New(db *gorm.DB) Storage {
	return Storage{
		DB: db,
	}
}

// Create inserts a new record for the given Table.
func (s Storage) Create(ctx context.Context, model table.Table) (uuid.UUID, error) {
	d := domainToStorage(model)

	res := s.DB.WithContext(ctx).Create(&d)
	if res.Error != nil {
		return uuid.Nil, res.Error
	}

	return d.ID, nil
}

// List returns all tables from the database.
func (s Storage) List(ctx context.Context) ([]table.Table, error) {
	var tables []Table

	res := s.DB.WithContext(ctx).Find(&tables)
	if res.Error != nil {
		return []table.Table{}, res.Error
	}

	return storageListToDomain(tables), nil
}
