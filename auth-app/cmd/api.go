package cmd

import (
	"auth-app/internal/api"
	"auth-app/internal/config"

	"github.com/spf13/cobra"
)

func apiCmd() *cobra.Command {
	cfg := config.Load()

	var command = &cobra.Command{
		Use:   "api",
		Short: "Run API server",
		Run: func(cmd *cobra.Command, args []string) {
			port := cfg.App.Port

			srv := api.NewServer()
			srv.Run(port)
		},
	}

	return command
}
