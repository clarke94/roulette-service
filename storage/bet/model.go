package bet

import (
	"time"

	"github.com/clarke94/roulette-service/internal/pkg/bet"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Bet is a storage model.
type Bet struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	TableID   uuid.UUID
	Bet       string
	Type      string
	Amount    int64
	Currency  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func domainToStorage(t *bet.Bet) Bet {
	return Bet{
		ID:       t.ID,
		TableID:  t.TableID,
		Bet:      t.Bet,
		Type:     t.Type,
		Amount:   t.Amount,
		Currency: t.Currency,
	}
}

func storageToDomain(t *Bet) bet.Bet {
	return bet.Bet{
		ID:       t.ID,
		TableID:  t.TableID,
		Bet:      t.Bet,
		Type:     t.Type,
		Amount:   t.Amount,
		Currency: t.Currency,
	}
}

func storageListToDomain(t []Bet) []bet.Bet {
	bets := make([]bet.Bet, len(t))

	for i := range t {
		bets[i] = storageToDomain(&t[i])
	}

	return bets
}

func domainListToStorage(t []bet.Bet) []Bet {
	bets := make([]Bet, len(t))

	for i := range t {
		bets[i] = domainToStorage(&t[i])
	}

	return bets
}
