package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Telegram BotConfig
	DB       DBConfig
}

type BotConfig struct {
	Token string
}

type DBConfig struct {
	Path string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	token := os.Getenv("TOKEN")
	if token == "" {
		return nil, fmt.Errorf("missing required env variable: TOKEN")
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "storage/suai_queue.db"
	}

	cfg := &Config{
		Telegram: BotConfig{Token: token},
		DB:       DBConfig{Path: dbPath},
	}

	return cfg, nil
}

func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(fmt.Sprintf("config load failed: %v", err))
	}
	return cfg
}
