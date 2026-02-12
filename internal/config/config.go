package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Conf struct {
	Server     HttpServerConf `yaml:"http-server"`
	DataSource DataSourceConf `yaml:"data-source"`
}

type HttpServerConf struct {
	Host           string        `yaml:"host"`
	Port           int           `yaml:"port"`
	RequestTimeout time.Duration `yaml:"request-timeout"`
	SessionTimeout time.Duration `yaml:"session-timeout"`
}

type DataSourceConf struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db-name"`
	DBHost   string `yaml:"host"`
	DBPort   int    `yaml:"port"`
	Driver   string `yaml:"driver"`
}

func MustLoad() Conf {
	configPath := os.Getenv("CONFIG")

	if configPath == "" {
		panic("Config path is empty")
	}

	var cfg Conf

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("Config couldn't read")
	}

	return cfg
}
