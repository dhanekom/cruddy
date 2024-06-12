package database

import (
	"context"

	"github.com/dhanekom/cruddy/internal/entities"
	"github.com/dhanekom/cruddy/internal/storage/cache"
)

type MockDBRepo struct {
	DBTables *cache.Cache[string, entities.DBTable]
}

func (r MockDBRepo) GetTables(ctx context.Context, schema string) ([]entities.DBTable, error) {
	tables := []entities.DBTable{}

	for _, table := range r.DBTables.Items() {
		tables = append(tables, entities.DBTable{
			Schema:        table.Schema,
			Tablename:     table.Tablename,
			FullTablename: table.FullTablename,
			Columns:       table.Columns,
		})
	}

	return tables, nil
}

func (r MockDBRepo) GetTableInfo(ctx context.Context, tablename string) (*entities.DBTable, error) {
	dbTable := entities.DBTable{}

	return &dbTable, nil
}
