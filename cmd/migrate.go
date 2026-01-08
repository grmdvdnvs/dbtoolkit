package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"grmdvdnvs/dbtoolkit/internal/logging"
	"grmdvdnvs/dbtoolkit/internal/migrate"
)

var (
	sourceDSN  string
	targetDSN  string
	tables     string
	strategy   string
	dryRun     bool
	retries    int
	waitBefore time.Duration
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrar datos entre bases heterogéneas",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := logging.Get()
		ctx := context.Background()

		if sourceDSN == "" || targetDSN == "" {
			return fmt.Errorf("se requieren --source y --target")
		}

		tableList := []string{}
		if tables != "" {
			for _, t := range strings.Split(tables, ",") {
				tableList = append(tableList, strings.TrimSpace(t))
			}
		}

		cfg := migrate.Config{
			SourceDSN: sourceDSN,
			TargetDSN: targetDSN,
			Tables:    tableList,
			Strategy:  strategy,
			DryRun:    dryRun,
			Retries:   retries,
			Backoff:   waitBefore,
		}

		planner := migrate.NewPlanner()
		plan, err := planner.Plan(ctx, cfg)
		if err != nil {
			logger.Error("plan failed", "err", err)
			return err
		}

		validator := migrate.NewValidator()
		if err := validator.ValidatePlan(ctx, plan); err != nil {
			logger.Error("validation failed", "err", err)
			return err
		}

		exec := migrate.NewExecutor()
		if err := exec.Execute(ctx, plan); err != nil {
			logger.Error("execution failed", "err", err)
			return err
		}

		logger.Info("migrate completed")
		return nil
	},
}

func init() {
	migrateCmd.Flags().StringVar(&sourceDSN, "source", "", "Source DSN (e.g. oracle://user:pass@host:port/sid)")
	migrateCmd.Flags().StringVar(&targetDSN, "target", "", "Target DSN (e.g. postgres://user:pass@host:port/db)")
	migrateCmd.Flags().StringVar(&tables, "tables", "", "Comma separated list of tables to migrate")
	migrateCmd.Flags().StringVar(&strategy, "strategy", "full", "Migration strategy: full|incremental")
	migrateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "No changes, only plan and validate")
	migrateCmd.Flags().IntVar(&retries, "retries", 3, "Number of retries on transient errors")
	migrateCmd.Flags().DurationVar(&waitBefore, "wait", 2*time.Second, "Initial backoff between retries")
}
