package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func connectToMysqlDB(dbConfigs DBConfigs) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbConfigs.Username, dbConfigs.Password, dbConfigs.Host, dbConfigs.Port, dbConfigs.DBName)

	db, err := sqlx.Open(dbConfigs.DriverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to DB: %s", err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("unable to ping db: %s", err)
	}

	return db, nil
}
