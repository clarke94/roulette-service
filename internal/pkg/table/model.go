package table

import "github.com/google/uuid"

// Table is a domain model.
type Table struct {
	ID         uuid.UUID
	Name       string `validate:"required"`
	MaximumBet int    `validate:"required,gte=10,gtefield=MinimumBet"`
	MinimumBet int    `validate:"required,gte=10"`
	Currency   string `validate:"required,oneof=GBP USD EUR"`
}
