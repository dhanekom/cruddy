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

func TestListTables(t *testing.T) {
	tests := []struct {
		name       string
		schema     string
		tableNames []string
		want       string
	}{
		{name: "no data", schema: "cruddb", tableNames: []string{}, want: ""},
		{name: "with data", schema: "cruddb", tableNames: []string{"sku\nstock"}, want: "sku\nstock"},
		{name: "invalid schema", schema: "", tableNames: []string{"sku\nstock"}, want: errorInvalidSchema.Error()},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			appConfigs := getAppConfigs(SetupData(test.schema, test.tableNames))

			crudGeneratorService := NewCRUDGeneratorService(appConfigs)
			buf := &bytes.Buffer{}
			err := crudGeneratorService.ListTables(context.Background(), buf, test.schema)
			if err != nil {
				if err.Error() != test.want {
					t.Fatalf("got error %q wanted %q", err.Error(), test.want)
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
			Schema:        schema,
			Tablename:     tablename,
			FullTablename: fmt.Sprintf("%s.%s", schema, tablename),
			Columns:       []entities.DBTableColumn{},
		})
	}

	return data
}
