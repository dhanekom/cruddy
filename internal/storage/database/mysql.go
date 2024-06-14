package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/dhanekom/cruddy/internal/entities"
	"github.com/jmoiron/sqlx"
)

func NewMysqlDBRepo(dbDriver *sqlx.DB) *MysqlDBRepo {
	return &MysqlDBRepo{
		db: dbDriver,
	}
}

type MysqlDBRepo struct {
	db *sqlx.DB
}

func (r MysqlDBRepo) GetTables(ctx context.Context, schema string) ([]entities.DBTable, error) {
	tables := []entities.DBTable{}

	err := r.db.SelectContext(ctx, &tables, `SELECT TABLE_SCHEMA, TABLE_NAME
	FROM INFORMATION_SCHEMA.TABLES
	WHERE TABLE_TYPE = 'BASE TABLE'
		AND TABLE_SCHEMA = ?;`, schema)

	if err != nil {
		return tables, err
	}

	if len(tables) == 0 {
		return tables, fmt.Errorf("no tables found for schema %q", schema)
	}

	return tables, nil
}

func (r MysqlDBRepo) GetTableInfo(ctx context.Context, schema, tablename string) (*entities.DBTable, error) {
	dbTable := entities.DBTable{}

	err := r.db.SelectContext(ctx, &dbTable.Columns, `SELECT c.column_name, c.data_type, 
	if(c.is_nullable = 'YES', true, false) IS_NULLABLE, if(kcu.CONSTRAINT_NAME = 'PRIMARY', true, false) IS_PK
FROM information_schema.columns c
LEFT OUTER JOIN information_schema.KEY_COLUMN_USAGE kcu on 
  kcu.TABLE_SCHEMA = c.TABLE_SCHEMA
  and kcu.TABLE_NAME = c.TABLE_NAME
  and kcu.COLUMN_NAME = c.COLUMN_NAME
WHERE (c.TABLE_SCHEMA = ?)
  AND (c.table_name = ?)`, schema, tablename)

	if err != nil {
		return &dbTable, err
	}

	if len(dbTable.Columns) == 0 {
		return &dbTable, fmt.Errorf("table %q not found", tablename)
	}

	tableNameParts := strings.Split(tablename, ".")
	if len(tableNameParts) == 2 {
		dbTable.Schema, dbTable.Tablename = tableNameParts[0], tableNameParts[1]
	} else {
		dbTable.Schema, dbTable.Tablename = "", tablename
	}

	return &dbTable, nil
}
