package sessions

import (
	"context"
	"time"

	"grmdvdnvs/dbtoolkit/internal/logging"
)

type KillConfig struct {
	DBType    string
	OlderThan time.Duration
	User      string
	Preview   bool
}

type Killer struct{}

func NewKiller() *Killer { return &Killer{} }

func (k *Killer) Kill(ctx context.Context, cfg KillConfig) error {
	logger := logging.Get()
	logger.Infof("kill sessions db=%s preview=%v", cfg.DBType, cfg.Preview)

	l := NewLister()
	sessions, err := l.List(ctx, ListConfig{DBType: cfg.DBType, OlderThan: cfg.OlderThan, User: cfg.User})
	if err != nil {
		return err
	}

	for _, s := range sessions {
		if cfg.Preview {
			logger.Infof("[preview] would kill session %s user=%s", s.ID, s.User)
			continue
		}
		logger.Infof("killing session %s user=%s", s.ID, s.User)
	}
	return nil
}
