package helpers

import (
	"github.com/jmoiron/sqlx"
)

func NewSqlConnection(connectionString string, logger Log) *sqlx.DB {

	if logger == nil {
		logger = NewConsoleLogger()
	}

	db, err := sqlx.Connect("mysql", connectionString)
	PanicToLogIf(err, logger)
	return db
}
