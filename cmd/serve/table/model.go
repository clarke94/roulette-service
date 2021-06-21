package table

import (
	"github.com/clarke94/roulette-service/internal/pkg/table"
)

// IDParam is the URL parameter binding the bet ID.
type IDParam struct {
	Table string `uri:"table" binding:"required,uuid"`
}

// Table is a presentation API model.
type Table struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name" binding:"required"`
	MaximumBet int    `json:"maximumBet" binding:"required,gte=10,gtefield=MinimumBet"`
	MinimumBet int    `json:"minimumBet" binding:"required,gte=10"`
	Currency   string `json:"currency" binding:"required,oneof=GBP USD EUR"`
}

// Update is a Table with a required ID binding.
type Update struct {
	ID string `json:"id" binding:"required,uuid"`
	Table
}

// Upsert is a presentation API model for the Upsert response.
type Upsert struct {
	ID string `json:"id"`
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
