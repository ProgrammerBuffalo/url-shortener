package main

import (
	"errors"
	"flag"
	"log/slog"
	"os"

	"github.com/ProgrammerBuffalo/url-shortener/internal/config"
	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	cfgPath := flag.String("config", "", "Path to YAML config file")
	flag.Parse()

	if *cfgPath == "" {
		logger.Error("Please provide --config flag")
		os.Exit(1)
	}

	args := flag.Args()

	if len(args) < 1 {
		logger.Error("Please provide command: up or down")
		os.Exit(1)
	}

	logger.Debug("Get configs")
	cfg := config.LoadMigrationConfig(*cfgPath)

	logger.Debug("Get migration file and connect to db")
	m, err := migrate.New(cfg.MigrationPath.Value, cfg.DataSource.Url)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	switch args[0] {
	case "up":
		logger.Debug("Up migration")
		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				logger.Debug("No migrations to apply")
				return
			}
			logger.Error(err.Error())
			os.Exit(1)
		}
	case "down":
		logger.Debug("Down migration")
		if err := m.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				logger.Debug("No migrations to apply")
				return
			}
			logger.Error(err.Error())
			os.Exit(1)
		}
	}

}
