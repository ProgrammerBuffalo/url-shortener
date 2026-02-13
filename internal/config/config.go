package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	Server     HttpServerConfig `yaml:"http-server"`
	DataSource DataSourceConfig `yaml:"data-source"`
}

type HttpServerConfig struct {
	Host           string        `yaml:"host"`
	Port           int           `yaml:"port"`
	RequestTimeout time.Duration `yaml:"request-timeout"`
	SessionTimeout time.Duration `yaml:"session-timeout"`
}

type DataSourceConfig struct {
	Url string `yaml:"url"`
}

type MigratorConfig struct {
	DataSource    DataSourceConfig `yaml:"data-source"`
	MigrationPath MigrationPath    `yaml:"migrations"`
}

type MigrationPath struct {
	Value string `yaml:"path"`
}

func LoadAppConfig(cfgPath string) *AppConfig {
	var cfg AppConfig
	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		log.Fatal(err.Error())
	}
	return &cfg
}

func LoadMigrationConfig(cfgPath string) *MigratorConfig {
	var cfg MigratorConfig
	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		log.Fatal(err.Error())
	}
	return &cfg
}
