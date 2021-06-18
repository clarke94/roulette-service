package table

// Table is a domain model.
type Table struct {
	Name       string `validate:"required"`
	MaximumBet int    `validate:"required,gte=10,gtefield=MinimumBet"`
	MinimumBet int    `validate:"required,gte=10"`
	Currency   string `validate:"required,oneof=GBP USD EUR"`
}
