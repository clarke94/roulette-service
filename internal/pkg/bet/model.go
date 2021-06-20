package bet

import (
	"github.com/google/uuid"
)

// Bet is a domain model.
type Bet struct {
	ID       uuid.UUID
	TableID  uuid.UUID `validate:"required"`
	Bet      string    `validate:"required"`
	Type     string    `validate:"required"`
	Amount   int64     `validate:"required"`
	Currency string    `validate:"required"`
}

// Update is a domain model for a bet update.
type Update struct {
	ID uuid.UUID `validate:"required"`
	Bet
}

// Type is the supported Bet type.
var (
	TypeRedBlack = "red/black"
	TypeOddEven  = "odd/even"
	TypeHighLow  = "high/low"
	TypeColumn   = "column"
	TypeDozen    = "dozen"
	TypeStraight = "straight"
)

// TypeMultiplierMap is the Bet type that is available and the associated multiplier for that bet.
var TypeMultiplierMap = map[string]int64{
	// Outside bets
	TypeRedBlack: 1,
	TypeOddEven:  1,
	TypeHighLow:  1,
	TypeColumn:   2,
	TypeDozen:    2,
	// Inside bets
	TypeStraight: 35,
}
