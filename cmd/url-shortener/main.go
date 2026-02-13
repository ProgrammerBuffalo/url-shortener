package main

import (
	"flag"
	"log/slog"
	"os"

	_ "github.com/lib/pq"

	"github.com/ProgrammerBuffalo/url-shortener/internal/config"
	"github.com/ProgrammerBuffalo/url-shortener/internal/storage"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	logger.Debug("Load config")

	cfgPath := flag.String("config", "", "Path to YAML config file")
	flag.Parse()

	if *cfgPath == "" {
		logger.Error("Please provide --config flag")
		os.Exit(1)
	}

	cfg := config.LoadAppConfig(*cfgPath)

	logger.Debug("Load database")

	if _, err := storage.New(cfg.DataSource); err != nil {
		logger.Error("Database load failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

}
