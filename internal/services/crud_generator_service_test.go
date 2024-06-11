package services

import (
	"bytes"
	"context"
	"slices"
	"strings"
	"testing"

	"github.com/dhanekom/cruddy/internal/configs"
	"github.com/dhanekom/cruddy/internal/entities"
	"github.com/dhanekom/cruddy/internal/storage/cache"
	"github.com/dhanekom/cruddy/internal/storage/database"
)

// var templatesFolder embed.FS

func TestGenerage(t *testing.T) {
	t.Run("valid scenario", func(t *testing.T) {

	})

	t.Run("invalid schema", func(t *testing.T) {

	})
}

func TestListTables(t *testing.T) {
	data := cache.New[string, entities.DBTable]()
	data.Set("db.stock", entities.DBTable{
		Schema:        "db",
		Tablename:     "stock",
		FullTablename: "db.stock",
		Columns:       []entities.DBTableColumn{},
	})
	data.Set("db.sku", entities.DBTable{
		Schema:        "db",
		Tablename:     "sku",
		FullTablename: "db.sku",
		Columns:       []entities.DBTableColumn{},
	})

	db := database.NewMockDBRepo(data)
	appConfigs := &configs.App{
		DB: db,
		// TemplatesFolder: fstest.MapFS{
		// 	//		"templates/crud_sql.templ": &fstest.MapFile{},
		// },
	}

	crudService := NewCRUDGeneratorService(appConfigs)

	buffer := &bytes.Buffer{}
	_ = crudService.ListTables(context.Background(), buffer, "db")
	result := strings.Split(buffer.String(), "\n")
	slices.Sort(result)
	got := strings.Join(result, "\n")
	want := `sku
stock`
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
