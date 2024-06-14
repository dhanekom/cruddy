package services

import (
	"bytes"
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/dhanekom/cruddy/internal/configs"
	"github.com/dhanekom/cruddy/internal/entities"
	"github.com/dhanekom/cruddy/internal/storage"
	"github.com/dhanekom/cruddy/internal/storage/cache"
	"github.com/dhanekom/cruddy/internal/storage/database"
)

const SQLC_STOCK_TEMPLATE = `-- name: GetPerson :one
select 
    id
  , first_name
  , last_name
  , email
from cruddb.person
where
  id = ?
limit 1;

-- name: ListPerson :many
select *
from cruddb.person;

-- name: CreatePerson :execresult
insert into cruddb.person (
    id
  , first_name
  , last_name
  , email
) values (
    ?
  , ?
  , ?
  , ?
);

-- name: UpdatePerson :exec
update cruddb.person
set
    id = ?
  , first_name = ?
  , last_name = ?
  , email = ?
where
  id = ?;

-- name: DeletePerson :exec
delete from cruddb.person
where
  id = ?;

`

func TestListTables(t *testing.T) {
	tests := []struct {
		name       string
		schema     string
		tableNames []string
		want       string
	}{
		{name: "no data", schema: "cruddb", tableNames: []string{}, want: ""},
		{name: "with data", schema: "cruddb", tableNames: []string{"sku", "stock"}, want: "cruddb.sku\ncruddb.stock"},
		{name: "invalid schema", schema: "", tableNames: []string{"sku", "stock"}, want: errorInvalidSchema.Error()},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			appConfigs := getAppConfigs(SetupData(test.schema, test.tableNames))

			crudGeneratorService := NewCRUDGeneratorService(appConfigs)
			buf := &bytes.Buffer{}
			err := crudGeneratorService.ListTables(context.Background(), buf, test.schema)
			if err != nil {
				if err.Error() != test.want {
					t.Fatalf("got error\n%q wanted\n%q", err.Error(), test.want)
				}
				return
			}

			tables := strings.Split(buf.String(), "\n")
			slices.Sort(tables)
			got := strings.Join(tables, "\n")

			if got != test.want {
				t.Errorf("got %q want %q", got, test.want)
			}
		})
	}
}

func TestGenerateSQLCTemplate(t *testing.T) {
	tests := []struct {
		name  string
		table entities.DBTable
		want  string
	}{
		{name: "valid table", table: entities.DBTable{Schema: "cruddb", Tablename: "person", Columns: []entities.DBTableColumn{
			{ColumnName: "id", IsPK: true},
			{ColumnName: "first_name"},
			{ColumnName: "last_name"},
			{ColumnName: "email"},
		}}, want: SQLC_STOCK_TEMPLATE},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data := cache.New[string, entities.DBTable]()
			data.Set(test.table.FullTablename(), test.table)
			appConfigs := getAppConfigs(data)

			crudGeneratorService := NewCRUDGeneratorService(appConfigs)
			buf := &bytes.Buffer{}
			err := crudGeneratorService.GenerateSQLCTemplate(context.Background(), buf, test.table.FullTablename())
			if err != nil {
				if err.Error() != test.want {
					t.Fatalf("got error %q wanted %q", err.Error(), test.want)
				}
				return
			}

			if buf.String() != test.want {
				t.Errorf("got\n%q\nwant\n%q", buf.String(), test.want)
			}
		})
	}
}

func TestTitle(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "empty string", input: "", want: ""},
		{name: "string with not spaces", input: "something", want: "Something"},
		{name: "string with spaces", input: "let's do something", want: "Let's do something"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := title(test.input)
			if got != test.want {
				t.Errorf("got %q want %q", got, test.want)
			}
		})
	}
}

func getAppConfigs(data *cache.Cache[string, entities.DBTable]) *configs.App {
	var dbRepo storage.DBRepo = database.MockDBRepo{
		DBTables: data,
	}
	return &configs.App{
		DB: dbRepo,
	}
}

func SetupData(schema string, tableNames []string) *cache.Cache[string, entities.DBTable] {
	data := cache.New[string, entities.DBTable]()

	for _, tablename := range tableNames {
		data.Set(fmt.Sprintf("%s.%s", schema, tablename), entities.DBTable{
			Schema:    schema,
			Tablename: tablename,
			Columns:   []entities.DBTableColumn{},
		})
	}

	return data
}
