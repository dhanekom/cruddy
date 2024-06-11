package services

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/dhanekom/cruddy/internal/configs"
)

//go:embed templates/crud_sql.templ
var crud_sql_template string

type CRUDGeneratorService struct {
	app *configs.App
}

func NewCRUDGeneratorService(a *configs.App) *CRUDGeneratorService {
	return &CRUDGeneratorService{
		app: a,
	}
}

func (s *CRUDGeneratorService) ListTables(ctx context.Context, w io.Writer, schema string) error {
	tables, err := s.app.DB.GetTables(ctx, schema)
	if err != nil {
		return err
	}

	var newline = ""
	for _, table := range tables {
		_, err = fmt.Fprintf(w, "%s%s", newline, table.Tablename)
		if err != nil {
			return err
		}
		newline = "\n"
	}

	return nil
}

func (s *CRUDGeneratorService) Generate(tablename string) error {
	tableInfo, err := s.app.DB.GetTableInfo(context.Background(), tablename)
	if err != nil {
		return err
	}

	templ, err := template.New("curd_sql").Parse(crud_sql_template)
	if err != nil {
		return err
	}

	err = templ.Execute(os.Stdout, tableInfo)
	if err != nil {
		return err
	}

	return nil
}
