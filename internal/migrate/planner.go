package migrate

import (
	"context"
	"time"

	"grmdvdnvs/dbtoolkit/internal/logging"
)

type Config struct {
	SourceDSN string
	TargetDSN string
	Tables    []string
	Strategy  string
	DryRun    bool
	Retries   int
	Backoff   time.Duration
}

type Plan struct {
	Config Config
	Steps  []Step
}

type Step struct {
	Table string
	Mode  string // full | incremental
}

type Planner struct{}

func NewPlanner() *Planner { return &Planner{} }

func (p *Planner) Plan(ctx context.Context, cfg Config) (Plan, error) {
	logger := logging.Get()
	logger.Infof("planning migration from %s to %s; dry-run=%v", cfg.SourceDSN, cfg.TargetDSN, cfg.DryRun)

	steps := []Step{}
	for _, t := range cfg.Tables {
		mode := "full"
		if cfg.Strategy == "incremental" {
			mode = "incremental"
		}
		steps = append(steps, Step{Table: t, Mode: mode})
	}

	plan := Plan{
		Config: cfg,
		Steps:  steps,
	}
	return plan, nil
}
