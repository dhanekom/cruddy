package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/dhanekom/cruddy/internal/configs"
	"github.com/dhanekom/cruddy/internal/services"
	"github.com/dhanekom/cruddy/internal/storage"
	"github.com/dhanekom/cruddy/internal/storage/database"
	"github.com/jmoiron/sqlx"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/urfave/cli/v2"
)

var db *sqlx.DB

var appConfigs *configs.App

func main() {
	// read configs
	//tablename := flag.String("t", "", "tablename for which to generate crud methods (format = 'schema.tablename')")

	config := koanf.New(".")

	var logLevel slog.Level
	if config.String("log.level") == "DEBUG" {
		logLevel = slog.LevelDebug
	} else {
		logLevel = slog.LevelInfo
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	if err := config.Load(file.Provider("configs.json"), json.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	dbConfigs := database.DBConfigs{
		DriverName: config.String("db.drivername"),
		Username:   config.String("db.username"),
		Password:   config.String("db.password"),
		Host:       config.String("db.host"),
		Port:       config.String("db.port"),
		DBName:     config.String("db.dbname"),
	}

	// connect to db driver
	var err error
	db, err = database.ConnectToDB(dbConfigs)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	slog.Debug("connected to db", slog.String("host", dbConfigs.Host), slog.String("port", dbConfigs.Port), slog.String("dbname", dbConfigs.DBName), slog.String("username", dbConfigs.Username))

	// Create mysql db repo
	var dbRepo storage.DBRepo
	dbRepo = database.NewMysqlDBRepo(db)

	appConfigs = &configs.App{
		DB: dbRepo,
		// Config:          config,
		// TemplatesFolder: templatesFolder,
	}

	app := &cli.App{
		Usage: "Generate CRUD code for a database table",
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "List tables in a schema",
				Subcommands: []*cli.Command{
					{
						Name:    "schema",
						Aliases: []string{"s"},
						Usage:   "Database schema name",
						Action: func(cCtx *cli.Context) error {
							crudService := services.NewCRUDGeneratorService(appConfigs)

							err = crudService.ListTables(context.Background(), os.Stdout, cCtx.Args().First())
							if err != nil {
								slog.Error(err.Error())
								os.Exit(1)
							}

							return nil
						},
					},
				},
			},
			{
				Name:    "generate",
				Aliases: []string{"g"},
				Usage:   "Generate CRUD code for a database table",
				Subcommands: []*cli.Command{
					{
						Name:    "tablename",
						Aliases: []string{"t"},
						Usage:   "DB table to generate code for (format = 'schema.tablename')",
						Action: func(cCtx *cli.Context) error {
							crudService := services.NewCRUDGeneratorService(appConfigs)

							err = crudService.GenerateSQLCTemplate(context.Background(), os.Stdout, cCtx.Args().First())
							if err != nil {
								slog.Error(err.Error())
								os.Exit(1)
							}

							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
