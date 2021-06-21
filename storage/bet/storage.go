package bet

import (
	"context"
	"errors"

	"github.com/clarke94/roulette-service/internal/pkg/bet"
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

// Create inserts a new record for the given Bet.
func (s Storage) Create(ctx context.Context, model bet.Bet) (string, error) {
	d := domainToStorage(&model)

	res := s.DB.WithContext(ctx).Create(&d)
	if res.Error != nil {
		return "", res.Error
	}

	return d.ID, nil
}

// List returns all bets from the database for a given table.
func (s Storage) List(ctx context.Context, tableID string, filters ...bet.Bet) ([]bet.Bet, error) {
	var bets []Bet

	queryFilters := domainListToStorage(filters)

	db := s.DB.WithContext(ctx)

	queryBuilder := db
	for _, filter := range queryFilters {
		queryBuilder = queryBuilder.Or(filter)
	}

	db = db.Where(&Bet{TableID: tableID}, queryBuilder)

	res := db.Find(&bets)

	if res.Error != nil {
		return []bet.Bet{}, res.Error
	}

	return storageListToDomain(bets), nil
}

// Update inserts a new record for the given Bet.
func (s Storage) Update(ctx context.Context, model bet.Bet) (string, error) {
	d := domainToStorage(&model)

	res := s.DB.WithContext(ctx).Model(&d).Updates(&d)
	if res.Error != nil {
		return "", res.Error
	}

	if res.RowsAffected == 0 {
		return "", errNoChange
	}

	return d.ID, nil
}

// Delete deletes a bet for the given table and ID.
func (s Storage) Delete(ctx context.Context, tableID, id string) (string, error) {
	res := s.DB.WithContext(ctx).Where(&Bet{TableID: tableID}).Delete(&Bet{}, id)
	if res.Error != nil {
		return "", res.Error
	}

	if res.RowsAffected == 0 {
		return "", errNoChange
	}

	return id, nil
}
