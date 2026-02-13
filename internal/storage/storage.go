package storage

import (
	"database/sql"
	"fmt"

	"github.com/ProgrammerBuffalo/url-shortener/internal/config"
)

type DataSource struct {
	db *sql.DB
}

func New(dsConf config.DataSourceConfig) (*DataSource, error) {
	db, err := sql.Open("postgres", dsConf.Url)

	if err != nil {
		return nil, fmt.Errorf("Couldnt open connection to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Couldnt ping to database: %w", err)
	}

	return &DataSource{db: db}, nil
}
