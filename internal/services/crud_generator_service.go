package services

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"strings"
	"text/template"
	"unicode"

	"github.com/dhanekom/cruddy/internal/configs"
)

//go:embed templates/crud_sql.templ
var crud_sql_template string

var errorInvalidSchema = errors.New("invalid schema")
var errorInvalidFullTablename = errors.New("invalid full tablename")

type CRUDGeneratorService struct {
	app *configs.App
}

func NewCRUDGeneratorService(a *configs.App) *CRUDGeneratorService {
	return &CRUDGeneratorService{
		app: a,
	}
}

func (s *CRUDGeneratorService) ListTables(ctx context.Context, out io.Writer, schema string) error {
	schema = strings.Trim(schema, " ")
	if schema == "" {
		return errorInvalidSchema
	}

	tables, err := s.app.DB.GetTables(ctx, schema)
	if err != nil {
		return fmt.Errorf("ListTables - unable to get tables: %w", err)
	}

	newline := ""
	for _, table := range tables {
		fmt.Fprintf(out, "%s%s", newline, table.FullTablename())
		newline = "\n"
	}

	return nil
}

func (s *CRUDGeneratorService) GenerateSQLCTemplate(ctx context.Context, out io.Writer, fullTablename string) error {
	fullTablename = strings.Trim(fullTablename, " ")
	if fullTablename == "" {
		return errorInvalidFullTablename
	}

	tablenameParts := strings.Split(fullTablename, ".")
	var schema, tablename string
	if len(tablenameParts) > 1 {
		schema = tablenameParts[0]
		tablename = tablenameParts[1]
	} else {
		schema = ""
		tablename = fullTablename
	}

	table, err := s.app.DB.GetTableInfo(ctx, schema, tablename)
	if err != nil {
		return fmt.Errorf("GenerateSQLCTemplate - GetTableInfo: %w", err)
	}

	funcMap := template.FuncMap{
		"title": title,
	}

	t, err := template.New("crud_sql").Funcs(funcMap).Parse(crud_sql_template)
	if err != nil {
		return fmt.Errorf("GenerateSQLCTemplate - template parsing: %w", err)
	}

	err = t.Execute(out, table)
	if err != nil {
		return fmt.Errorf("GenerateSQLCTemplate - generate template: %w", err)
	}

	return nil
}

func title(input string) string {
	var sb strings.Builder
	for i, r := range input {
		if i == 0 {
			sb.WriteString(string(unicode.ToTitle(r)))
		} else {
			sb.WriteString(string(r))
		}
	}

	return sb.String()
}
