package bet

import (
	"github.com/clarke94/roulette-service/internal/pkg/bet"
	"github.com/google/uuid"
)

// Bet is a presentation API model.
type Bet struct {
	ID       uuid.UUID `json:"id,omitempty"`
	Bet      string    `json:"bet"`
	Type     string    `json:"type"`
	Amount   int64     `json:"amount"`
	Currency string    `json:"currency"`
}

// Upsert is a presentation API model for the Upsert response.
type Upsert struct {
	ID uuid.UUID `json:"id"`
}

// Error is a presentation API model for the Error response.
type Error struct {
	Error string `json:"error"`
}

func presentationToDomain(t Bet, tableID uuid.UUID) bet.Bet {
	return bet.Bet{
		ID:       t.ID,
		TableID:  tableID,
		Bet:      t.Bet,
		Type:     t.Type,
		Amount:   t.Amount,
		Currency: t.Currency,
	}
}

func domainToPresentation(t *bet.Bet) Bet {
	return Bet{
		ID:       t.ID,
		Bet:      t.Bet,
		Type:     t.Type,
		Amount:   t.Amount,
		Currency: t.Currency,
	}
}

func domainListToPresentation(t []bet.Bet) []Bet {
	bets := make([]Bet, len(t))

	for i := range t {
		bets[i] = domainToPresentation(&t[i])
	}

	return bets
}
