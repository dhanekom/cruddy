package entities

import "fmt"

type DBTable struct {
	Schema    string `db:"TABLE_SCHEMA"`
	Tablename string `db:"TABLE_NAME"`
	Columns   []DBTableColumn
}

type DBTableColumn struct {
	ColumnName string `db:"COLUMN_NAME"`
	IsPK       bool   `db:"IS_PK"`
	DataType   string `db:"DATA_TYPE"`
	IsNullable bool   `db:"IS_NULLABLE"`
}

func (t DBTable) FullTablename() string {
	if t.Schema != "" {
		return fmt.Sprintf("%s.%s", t.Schema, t.Tablename)
	} else {
		return t.Tablename
	}
}

// func (t DBTable)
