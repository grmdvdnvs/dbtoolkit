package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type PostgresConnector struct {
	DSN string
	DB  *sql.DB
}

func New(dsn string) (*PostgresConnector, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &PostgresConnector{DSN: dsn, DB: db}, nil
}
