package service

import (
	"context"
	"log/slog"

	"github.com/ProgrammerBuffalo/url-shortener/internal/errs"
	"github.com/ProgrammerBuffalo/url-shortener/internal/repository"
	"github.com/ProgrammerBuffalo/url-shortener/internal/repository/model"
	"github.com/ProgrammerBuffalo/url-shortener/internal/storage"
	"github.com/ProgrammerBuffalo/url-shortener/lib"
	"github.com/google/uuid"
)

type UrlService struct {
	logger *slog.Logger
	ds     *storage.DataSource
	r      repository.Repository
}

func NewUrlService(ds *storage.DataSource, logger *slog.Logger, r repository.Repository) *UrlService {
	return &UrlService{ds: ds, logger: logger, r: r}
}

func (s *UrlService) FindByURL(ctx context.Context, shortUrl string) (string, error) {
	shortUrl, err := s.r.Read(ctx, shortUrl)

	if err != nil {
		return "", err
	}

	return shortUrl, nil
}

func (s *UrlService) Create(ctx context.Context, longUrl string) (string, error) {
	s.logger.Info("Call create url service")
	id, err := uuid.NewRandom()

	if err != nil {
		return "", errs.ErrInternalServer
	}

	urlDao := model.UrlDao{
		Id:    id,
		Short: lib.Hash(id),
		Long:  longUrl,
	}

	s.logger.Info("Start transaction url service")
	tx, err := s.ds.SqlDb.BeginTx(ctx, nil)

	defer func() {
		if err := tx.Commit(); err != nil {
			s.logger.Error("Transaction commit failed url service", slog.String("error", err.Error()))
		}
	}()

	if err != nil {
		return "", errs.ErrInternalServer
	}

	s.logger.Info("Create resource url service")
	if err := s.r.Create(ctx, tx, &urlDao); err != nil {
		if err := tx.Rollback(); err != nil {
			s.logger.Error("Transaction rollback failed url service", slog.String("error", err.Error()))
		}
		return "", errs.ErrInternalServer
	}

	s.logger.Info("Resource url service created", slog.String("url", urlDao.Short))

	return urlDao.Short, nil
}
