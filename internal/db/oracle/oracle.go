package oracle

import (
	"database/sql"

	_ "github.com/godror/godror"
)

type OracleConnector struct {
	DSN string
	DB  *sql.DB
}

func New(dsn string) (*OracleConnector, error) {
	db, err := sql.Open("godror", dsn)
	if err != nil {
		return nil, err
	}
	return &OracleConnector{DSN: dsn, DB: db}, nil
}
