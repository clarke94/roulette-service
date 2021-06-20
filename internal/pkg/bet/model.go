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

// Result is the round result from a game.
type Result struct {
	Number  int      `json:"number"`
	Color   string   `json:"color"`
	Winners []Winner `json:"winners"`
}

// Winner is a winning bet from a round.
type Winner struct {
	BetID    uuid.UUID `json:"betId"`
	Amount   int64     `json:"amount"`
	Currency string    `json:"currency"`
}

const (
	colorRed   = "red"
	colorBlack = "black"
	colorGreen = "green"
)

// Type is the supported Bet type.
var (
	TypeRedBlack = "red/black"
	TypeStraight = "straight"
)

// TypeMultiplierMap is the Bet type that is available and the associated multiplier for that bet.
var TypeMultiplierMap = map[string]int64{
	// Outside bets
	TypeRedBlack: 1,
	// Inside bets
	TypeStraight: 35,
}

func betToWinner(b Bet) Winner {
	return Winner{
		BetID:    b.ID,
		Amount:   b.Amount,
		Currency: b.Currency,
	}
}

func betListToWinner(b []Bet) []Winner {
	winners := make([]Winner, len(b))

	for i := range b {
		winners[i] = betToWinner(b[i])
	}

	return winners
}
