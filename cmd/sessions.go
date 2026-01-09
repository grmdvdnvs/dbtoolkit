package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"grmdvdnvs/dbtoolkit/internal/logging"
	"grmdvdnvs/dbtoolkit/internal/sessions"
)

var (
	dsn        string
	driver     string
	olderThan  time.Duration
	userFilter string
	kill       bool
	preview    bool
)

var sessionsCmd = &cobra.Command{
	Use:   "sessions",
	Short: "Auditar y cerrar sesiones activas en bases de datos",
}

func init() {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Listar sesiones",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := logging.Get()
			ctx := context.Background()
			cfg := sessions.ListConfig{
				DSN:       dsn,
				Driver:    driver,
				OlderThan: olderThan,
				User:      userFilter,
			}
			lister := sessions.NewLister()
			res, err := lister.List(ctx, cfg)
			if err != nil {
				logger.Error("list sessions failed", "err", err)
				return err
			}
			for _, r := range res {
				fmt.Println(r.String())
			}
			return nil
		},
	}

	killCmd := &cobra.Command{
		Use:   "kill",
		Short: "Cerrar sesiones que cumplan filtros",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := logging.Get()
			ctx := context.Background()
			cfg := sessions.KillConfig{
				Driver:    driver,
				OlderThan: olderThan,
				User:      userFilter,
				Preview:   preview,
			}
			k := sessions.NewKiller()
			logger.Infof("starting kill sessions (preview=%v)", preview)
			return k.Kill(ctx, cfg)
		},
	}

	sessionsCmd.PersistentFlags().StringVar(&dsn, "dsn", "", "Data Source Name de la BD (e.g. user/password@host:port/sid)")
	sessionsCmd.PersistentFlags().StringVar(&driver, "driver", "", "Tipo de BD: oracle|postgres")
	sessionsCmd.PersistentFlags().DurationVar(&olderThan, "older-than", 0, "Listar/kill sesiones más antiguas que (e.g. 2h)")
	sessionsCmd.PersistentFlags().StringVar(&userFilter, "user", "", "Filtrar por usuario")
	listCmd.Flags().BoolVar(&preview, "preview", true, "Modo seguro: sólo ver qué se cerraría")
	killCmd.Flags().BoolVar(&preview, "preview", false, "Si false se realizará el cierre efectivo")

	sessionsCmd.AddCommand(listCmd)
	sessionsCmd.AddCommand(killCmd)
}
