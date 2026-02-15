package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ProgrammerBuffalo/url-shortener/internal/errs"
	"github.com/ProgrammerBuffalo/url-shortener/internal/repository/model"
	"github.com/ProgrammerBuffalo/url-shortener/internal/storage"
	"github.com/lib/pq"
)

const (
	pgErrUniqueViolation = "23505"
)

type Repository interface {
	CreateTx(ctx context.Context, tx *sql.Tx, dao *model.UrlDao) error
	ExistsTx(ctx context.Context, tx *sql.Tx, shortUrl string) (bool, error)
	Read(ctx context.Context, shortUrl string) (string, error)
}

type UrlRepository struct {
	ds *storage.DataSource
}

func NewUrlRepository(ds *storage.DataSource) *UrlRepository {
	return &UrlRepository{ds: ds}
}

func (r *UrlRepository) CreateTx(ctx context.Context, tx *sql.Tx, dao *model.UrlDao) error {
	_, err := tx.ExecContext(
		ctx,
		"INSERT INTO url(id, long_url, short_url) VALUES ($1, $2, $3)",
		dao.Id, dao.Long, dao.Short)

	if err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("db error: %w", errs.ErrDuplicateUrl)
		}
		return fmt.Errorf("db error: %w", err)
	}

	return nil
}

func (r *UrlRepository) Read(ctx context.Context, shortUrl string) (string, error) {
	var url string
	err := r.ds.SqlDb.QueryRowContext(ctx, "SELECT short_url FROM url WHERE short_url = $1", shortUrl).Scan(&url)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errs.ErrUrlNotFound
		}
		return "", fmt.Errorf("db error: %w", errs.ErrInternalServer)
	}

	return url, nil
}

func (r *UrlRepository) ExistsTx(ctx context.Context, tx *sql.Tx, longUrl string) (bool, error) {
	var exists int
	if err := tx.QueryRowContext(ctx, "SELECT 1 FROM url WHERE long_url = $1", longUrl).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("db error: %w", errs.ErrInternalServer)
	}
	return exists > 0, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgErrUniqueViolation {
			return true
		}
	}
	return false
}
