package main

import (
	"log/slog"
	"os"

	_ "github.com/lib/pq"

	"github.com/ProgrammerBuffalo/url-shortener/internal/config"
	"github.com/ProgrammerBuffalo/url-shortener/internal/storage"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	logger.Debug("Load config")

	cfg := config.MustLoad()

	logger.Debug("Load database")

	if _, err := storage.New(cfg.DataSource); err != nil {
		logger.Error("Database load failed", slog.String("error", err.Error()))
		panic("Database error")
	}

}
