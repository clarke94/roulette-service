package table

import (
	"context"
	"errors"

	"github.com/clarke94/roulette-service/internal/pkg/table"
	"gorm.io/gorm"
)

var errNoChange = errors.New("no change")

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
func (s Storage) Create(ctx context.Context, model table.Table) (string, error) {
	d := domainToStorage(model)

	res := s.DB.WithContext(ctx).Create(&d)
	if res.Error != nil {
		return "", res.Error
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

// Update inserts a new record for the given Table.
func (s Storage) Update(ctx context.Context, model table.Table) (string, error) {
	d := domainToStorage(model)

	res := s.DB.WithContext(ctx).Model(&d).Updates(&d)
	if res.Error != nil {
		return "", res.Error
	}

	if res.RowsAffected == 0 {
		return "", errNoChange
	}

	return d.ID, nil
}

// Delete deletes a table for the given ID.
func (s Storage) Delete(ctx context.Context, id string) (string, error) {
	res := s.DB.WithContext(ctx).Delete(&Table{}, id)
	if res.Error != nil {
		return "", res.Error
	}

	if res.RowsAffected == 0 {
		return "", errNoChange
	}

	return id, nil
}
