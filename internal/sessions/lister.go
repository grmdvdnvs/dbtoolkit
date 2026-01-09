package sessions

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"grmdvdnvs/dbtoolkit/internal/logging"
)

type SessionInfo struct {
	PID         int
	User        sql.NullString
	Database    sql.NullString
	Application sql.NullString
	ClientAddr  sql.NullString
	State       sql.NullString
	QueryStart  sql.NullTime
	Query       sql.NullString
}

func (s SessionInfo) String() string {
	return fmt.Sprintf(
		"PID: %d, User: %s, Database: %s, Application: %s, ClientAddr: %s, State: %s, QueryStart: %s, Query: %s",
		s.PID,
		ns(s.User),
		ns(s.Database),
		ns(s.Application),
		ns(s.ClientAddr),
		ns(s.State),
		nt(s.QueryStart),
		ns(s.Query),
	)
}

func ns(v sql.NullString) string {
	if v.Valid {
		return v.String
	}
	return "-"
}

func nt(v sql.NullTime) string {
	if v.Valid {
		return v.Time.Format(time.RFC3339)
	}
	return "-"
}

type ListConfig struct {
	DSN       string
	Driver    string
	OlderThan time.Duration
	User      string
}

type Lister struct{}

func NewLister() *Lister { return &Lister{} }

func (l *Lister) List(ctx context.Context, cfg ListConfig) ([]SessionInfo, error) {
	logger := logging.Get()
	logger.Infof("listing sessions for driver=%s dsn=%s", cfg.Driver, cfg.DSN)
	db, err := sql.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return listActiveSessions(db)
}

func listActiveSessions(db *sql.DB) ([]SessionInfo, error) {
	query := `
		SELECT
			pid,
			usename,
			datname,
			application_name,
			client_addr::text,
			state,
			query_start,
			query
		FROM pg_stat_activity
		WHERE pid <> pg_backend_pid()
		AND backend_type = 'client backend'
		ORDER BY state, query_start;
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []SessionInfo

	for rows.Next() {
		var s SessionInfo
		if err := rows.Scan(
			&s.PID,
			&s.User,
			&s.Database,
			&s.Application,
			&s.ClientAddr,
			&s.State,
			&s.QueryStart,
			&s.Query,
		); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}

	return sessions, rows.Err()
}
