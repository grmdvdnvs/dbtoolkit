package migrate

import (
	"context"
	"errors"

	"grmdvdnvs/dbtoolkit/internal/logging"
)

type Validator struct{}

func NewValidator() *Validator { return &Validator{} }

func (v *Validator) ValidatePlan(ctx context.Context, plan Plan) error {
	logger := logging.Get()
	if len(plan.Steps) == 0 {
		return errors.New("plan vacío: sin tablas a migrar")
	}
	for _, s := range plan.Steps {
		logger.Infof("validate step: table=%s mode=%s", s.Table, s.Mode)
	}
	return nil
}
