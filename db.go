package helpers

import (
	"github.com/jmoiron/sqlx"
	"log"
)

func NewSqlConnection(connectionString string, logger *log.Logger) *sqlx.DB {

	if logger == nil {
		logger = Logger()
	}

	db, err := sqlx.Connect("mysql", connectionString)
	PanicIf(err, logger)
	return db
}
