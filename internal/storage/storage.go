package storage

import (
	"database/sql"
	"fmt"

	"github.com/ProgrammerBuffalo/url-shortener/internal/config"
)

type DataSource struct {
	db *sql.DB
}

func New(dsConf config.DataSourceConf) (*DataSource, error) {
	db, err := sql.Open(dsConf.Driver, getConnectionStr(&dsConf))

	if err != nil {
		return nil, fmt.Errorf("Couldnt open connection to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Couldnt ping to database: %w", err)
	}

	return &DataSource{db: db}, nil
}

func getConnectionStr(dsConf *config.DataSourceConf) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dsConf.Username,
		dsConf.Password,
		dsConf.DBHost,
		dsConf.DBPort,
		dsConf.DBName,
	)
}
