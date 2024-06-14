package database

import (
	"context"
	"fmt"

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
			Schema:    table.Schema,
			Tablename: table.Tablename,
			Columns:   table.Columns,
		})
	}

	return tables, nil
}

func (r MockDBRepo) GetTableInfo(ctx context.Context, schema, tablename string) (*entities.DBTable, error) {
	dbTable := entities.DBTable{}

	var fullTablename string
	if schema != "" {
		fullTablename = fmt.Sprintf("%s.%s", schema, tablename)
	} else {
		fullTablename = tablename
	}

	dbTable, found := r.DBTables.Get(fullTablename)
	if !found {
		return &dbTable, fmt.Errorf("table %q not found", fullTablename)
	}

	return &dbTable, nil
}
