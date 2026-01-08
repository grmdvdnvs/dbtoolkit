// File: cmd/root.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"grmdvdnvs/dbtoolkit/internal/logging"
)

var rootCmd = &cobra.Command{
	Use:   "dbtoolkit",
	Short: "Herramienta CLI para tareas DB: migraciones y gestión de sesiones",
}

func Execute() {
	logging.InitDefault()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// flags comunes si los necesitas
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(sessionsCmd)
}
