package storage

import (
	"context"

	"github.com/dhanekom/cruddy/internal/entities"
)

type DBRepo interface {
	GetTables(ctx context.Context, schema string) ([]entities.DBTable, error)
	GetTableInfo(ctx context.Context, schema, tablename string) (*entities.DBTable, error)
}
