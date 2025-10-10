package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

var ErrCfgInvalid = errors.New("invalid configuration")

type Config struct {
	Env   string `env:"APP_ENV"`
	DB    DBConfig
	HTTP  HTTPConfig
	Cache CacheConfig
	Redis RedisConfig
}

type DBConfig struct {
	Host            string        `env:"DB_HOST"`
	Port            int           `env:"DB_PORT"`
	User            string        `env:"DB_USER"`
	Password        string        `env:"DB_PASSWORD"`
	Name            string        `env:"DB_NAME"`
	SSLMode         string        `env:"DB_SSL"`
	MaxOpenConns    int32         `env:"DB_MAX_CONNS"`
	MaxIdleConns    int32         `env:"DB_IDLE_CONNS"`
	ConnMaxLifetime time.Duration `env:"DB_CONN_TIME_LIFE"`
	ConnMaxIdleTime time.Duration `env:"DB_CONN_TIME_IDLE"`
}

type HTTPConfig struct {
	Host         string `env:"HTTP_HOST"`
	Port         string `env:"HTTP_PORT"`
	ReadTimeout  int    `env:"HTTP_READ_TIMEOUT" env-default:"5"`
	WriteTimeout int    `env:"HTTP_WRITE_TIMEOUT" env-default:"10"`
	IdleTimeout  int    `env:"HTTP_IDLE_TIMEOUT" env-default:"120"`
}

type CacheConfig struct {
	Limit int `env:"CACHE_LIMIT" env-default:"1000"`
}

type RedisConfig struct {
	Addr     string `env:"REDIS_ADDR"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB"`
	TTL      int    `env:"REDIS_TTL"`
}

func (dc *DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dc.Host, dc.Port, dc.User, dc.Password, dc.Name, dc.SSLMode,
	)
}

func InitConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	cfg, err := mapStructs()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func mapStructs() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
