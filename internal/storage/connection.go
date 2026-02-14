package storage

import (
	"database/sql"
	"fmt"

	"github.com/ProgrammerBuffalo/url-shortener/internal/config"
)

type DataSource struct {
	SqlDb *sql.DB
}

func New(dsConf config.DataSourceConfig) (*DataSource, error) {
	db, err := sql.Open("postgres", dsConf.Url)

	if err != nil {
		return nil, fmt.Errorf("couldn't open connection to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("couldn't ping to database: %w", err)
	}

	return &DataSource{SqlDb: db}, nil
}
