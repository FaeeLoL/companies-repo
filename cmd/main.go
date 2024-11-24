package main

import (
	"context"
	"errors"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/faeelol/companies-store/internal/app"
	"github.com/faeelol/companies-store/internal/app/database"
)

const (
	version = "1.0.0"
)

var configPath string

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "companies-repo",
		Short: "Companies repo",
		RunE: func(_ *cobra.Command, _ []string) error {
			return errors.New("choose setup mode")
		},
	}
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "configs/config.example.yml", "path to config file")

	rootCmd.AddCommand(NewHTTPServerCommand())
	rootCmd.AddCommand(NewMigrateDBCommand())

	rootCmd.Version = version
	return rootCmd
}

func NewHTTPServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "http",
		Short: "HTTP server",
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg := app.NewConfig()
			return app.LoadConfigInitLoggerAndDo(configPath, cfg, func(logger *logrus.Logger) error {
				return app.StartHTTPService(context.Background(), cfg, logger)
			})
		},
	}
	return cmd
}

func NewMigrateDBCommand() *cobra.Command {
	var migrateDBDown bool
	migrateDBCmd := &cobra.Command{
		Use:   "migrate-db",
		Short: "Migrate database",
		RunE: func(_ *cobra.Command, _ []string) error {
			direction := database.MigrDirectionUp
			if migrateDBDown {
				direction = database.MigrDirectionDown
			}
			cfg := app.NewConfig()
			return app.LoadConfigInitLoggerAndDo(configPath, cfg, func(logger *logrus.Logger) error {
				return app.MigrateDatabase(context.Background(), cfg, direction, logger)
			})
		},
	}
	migrateDBCmd.Flags().BoolVar(&migrateDBDown, "down", false, "migrate down")
	return migrateDBCmd
}

func main() {
	if err := NewRootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
