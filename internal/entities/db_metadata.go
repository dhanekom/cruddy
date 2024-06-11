package entities

type DBTable struct {
	Schema        string `db:"TABLE_SCHEMA"`
	Tablename     string `db:"TABLE_NAME"`
	FullTablename string
	Columns       []DBTableColumn
}

type DBTableColumn struct {
	ColumnName string `db:"COLUMN_NAME"`
	IsPK       bool   `db:"IS_PK"`
	DataType   string `db:"DATA_TYPE"`
	IsNullable bool   `db:"IS_NULLABLE"`
}
