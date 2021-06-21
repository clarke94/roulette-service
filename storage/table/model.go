package table

import (
	"time"

	"github.com/clarke94/roulette-service/internal/pkg/table"
	"gorm.io/gorm"
)

// Table is a storage model.
type Table struct {
	ID         string `gorm:"primaryKey"`
	Name       string
	MaximumBet int
	MinimumBet int
	Currency   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func domainToStorage(t table.Table) Table {
	return Table{
		ID:         t.ID,
		Name:       t.Name,
		MaximumBet: t.MaximumBet,
		MinimumBet: t.MinimumBet,
		Currency:   t.Currency,
	}
}

func storageToDomain(t *Table) table.Table {
	return table.Table{
		ID:         t.ID,
		Name:       t.Name,
		MaximumBet: t.MaximumBet,
		MinimumBet: t.MinimumBet,
		Currency:   t.Currency,
	}
}

func storageListToDomain(t []Table) []table.Table {
	tables := make([]table.Table, len(t))

	for i := range t {
		tables[i] = storageToDomain(&t[i])
	}

	return tables
}
