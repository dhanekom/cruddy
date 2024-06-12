package services

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io"

	"github.com/dhanekom/cruddy/internal/configs"
)

//go:embed templates/crud_sql.templ
var crud_sql_template string

var errorInvalidSchema = errors.New("valid schema required")

type CRUDGeneratorService struct {
	app *configs.App
}

func NewCRUDGeneratorService(a *configs.App) *CRUDGeneratorService {
	return &CRUDGeneratorService{
		app: a,
	}
}

func (s *CRUDGeneratorService) ListTables(ctx context.Context, out io.Writer, schema string) error {
	if schema == "" {
		return errorInvalidSchema
	}

	tables, err := s.app.DB.GetTables(ctx, schema)
	if err != nil {
		return fmt.Errorf("ListTables - unable to get tables: %w", err)
	}

	newline := ""
	for _, table := range tables {
		fmt.Fprintf(out, "%s%s", newline, table.Tablename)
		newline = "\n"
	}

	return nil
}
