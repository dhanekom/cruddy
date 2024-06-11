package database

import "github.com/jmoiron/sqlx"

type DBConfigs struct {
	DriverName string // e.g. pgx, mysql etc
	Username   string
	Password   string
	Host       string
	Port       string
	DBName     string
}

func ConnectToDB(dbConfigs DBConfigs) (*sqlx.DB, error) {
	return connectToMysqlDB(dbConfigs)
}
