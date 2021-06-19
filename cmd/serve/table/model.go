package table

import (
	"github.com/clarke94/roulette-service/internal/pkg/table"
	"github.com/google/uuid"
)

// Table is a presentation API model.
type Table struct {
	ID         uuid.UUID `json:"id,omitempty"`
	Name       string    `json:"name"`
	MaximumBet int       `json:"maximumBet"`
	MinimumBet int       `json:"minimumBet"`
	Currency   string    `json:"currency"`
}

// Create is a presentation API model for the Create response.
type Create struct {
	ID uuid.UUID `json:"id"`
}

// Error is a presentation API model for the Error response.
type Error struct {
	Error string `json:"error"`
}

func presentationToDomain(t Table) table.Table {
	return table.Table(t)
}

func domainToPresentation(t table.Table) Table {
	return Table(t)
}

func domainListToPresentation(t []table.Table) []Table {
	tables := make([]Table, len(t))

	for i := range t {
		tables[i] = domainToPresentation(t[i])
	}

	return tables
}
