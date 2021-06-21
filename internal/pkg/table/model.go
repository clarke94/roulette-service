package table

// Table is a domain model.
type Table struct {
	ID         string
	Name       string
	MaximumBet int
	MinimumBet int
	Currency   string
}
