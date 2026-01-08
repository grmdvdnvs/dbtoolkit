package migrate

import (
	"context"
	"time"

	"grmdvdnvs/dbtoolkit/internal/logging"
)

type Executor struct{}

func NewExecutor() *Executor { return &Executor{} }

func (e *Executor) Execute(ctx context.Context, plan Plan) error {
	logger := logging.Get()
	cfg := plan.Config

	for _, step := range plan.Steps {
		logger.Infof("execute step: table=%s mode=%s", step.Table, step.Mode)
		if cfg.DryRun {
			logger.Infof("[dry-run] would migrate table %s", step.Table)
			continue
		}

		attempts := 0
		for {
			attempts++
			err := mockCopy(ctx, step)
			if err == nil {
				logger.Infof("migrated table %s", step.Table)
				break
			}
			if attempts >= cfg.Retries {
				return err
			}
			logger.Warnf("retrying table %s after error: %v", step.Table, err)
			time.Sleep(cfg.Backoff)
		}
	}
	return nil
}

func mockCopy(ctx context.Context, s Step) error {
	// placeholder: implement actual row copy and transforms
	return nil
}
