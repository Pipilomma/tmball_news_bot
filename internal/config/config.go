package config

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env             string        `env:"ENV" env-default:"local"`
	GracefulTimeout time.Duration `env:"GRACEFUL_TIMEOUT" env-default:"10s"`
	Parser          ParserConfig
	DB              DBConfig
	Telegram        TelegramConfig
}

type ParserConfig struct {
	ParserTimeout time.Duration `env:"PARSER_TIMEOUT" env-default:"1h"`
	ParserUrl     string        `env:"URL_FOR_PARSER" required:"true"`
}

type DBConfig struct {
	DBName   string `env:"POSTGRES_DB" required:"true"`
	Username string `env:"POSTGRES_USER" required:"true"`
	Port     string `env:"POSTGRES_PORT" required:"true"`
	Host     string `env:"POSTGRES_HOST" required:"true"`
	SSLMode  string `env:"POSTGRES_SSLMODE" required:"true"`
	Password string `env:"POSTGRES_PASSWORD" required:"true"`
}

type TelegramConfig struct {
	Enabled    bool   `env:"TELEGRAM_ENABLED" env-default:"true"`
	BotToken   string `env:"TELEGRAM_BOT_TOKEN" env-required:"true"`
	AuthorName string `env:"TELEGRAM_BOT_AUTHOR_NAME"`
	Timeout    int    `env:"TELEGRAM_BOT_TIMEOUT" env-default:"60"`
	Debug      bool   `env:"TELEGRAM_BOT_DEBUG" env-default:"false"`
}

func MustLoad() *Config {
	var cfg Config

	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("No loading .env file: %v", err)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("No loading env variables: %v", err)
	}

	if cfg.Telegram.Enabled {
		if cfg.Telegram.BotToken == "" {
			log.Fatal("TELEGRAM_BOT_TOKEN is required when TELEGRAM_ENABLED=true")
		}
	}

	if !cfg.Telegram.Enabled {
		log.Fatal("At least one bot must be enabled")
	}

	return &cfg
}
