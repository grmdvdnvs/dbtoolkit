package sessions

import (
	"context"
	"fmt"
	"time"

	"grmdvdnvs/dbtoolkit/internal/logging"
)

type SessionInfo struct {
	ID        string
	User      string
	StartedAt time.Time
	State     string
}

func (s SessionInfo) String() string {
	return fmt.Sprintf("id=%s user=%s started=%s state=%s", s.ID, s.User, s.StartedAt.Format(time.RFC3339), s.State)
}

type ListConfig struct {
	DBType    string
	OlderThan time.Duration
	User      string
}

type Lister struct{}

func NewLister() *Lister { return &Lister{} }

func (l *Lister) List(ctx context.Context, cfg ListConfig) ([]SessionInfo, error) {
	logger := logging.Get()
	logger.Infof("listing sessions for db=%s older-than=%v user=%s", cfg.DBType, cfg.OlderThan, cfg.User)

	now := time.Now()
	sessions := []SessionInfo{
		{ID: "1", User: "app_user", StartedAt: now.Add(-3 * time.Hour), State: "ACTIVE"},
		{ID: "2", User: "other", StartedAt: now.Add(-30 * time.Minute), State: "IDLE"},
	}
	out := []SessionInfo{}
	for _, s := range sessions {
		if cfg.User != "" && s.User != cfg.User {
			continue
		}
		if cfg.OlderThan > 0 && now.Sub(s.StartedAt) < cfg.OlderThan {
			continue
		}
		out = append(out, s)
	}
	return out, nil
}
