package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ProgrammerBuffalo/url-shortener/internal/http/handler"
	"github.com/ProgrammerBuffalo/url-shortener/internal/repository"
	"github.com/ProgrammerBuffalo/url-shortener/internal/service"
	"github.com/go-chi/chi/v5"
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

	ds, err := storage.New(cfg.DataSource)
	if err != nil {
		logger.Error("Database load failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	r := repository.NewUrlRepository(ds)
	s := service.NewUrlService(ds, logger, r)

	router := chi.NewRouter()

	router.Post("/", handler.NewSaveHandler(s, logger).ServeHTTP)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.RequestTimeout,
		WriteTimeout: cfg.Server.RequestTimeout,
		IdleTimeout:  cfg.Server.SessionTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("failed to start server", err.Error())
		}
	}()

	logger.Info("Server started")

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("failed to shutdown server", err.Error())
		return
	}

	logger.Info("Server shutdown")
}
