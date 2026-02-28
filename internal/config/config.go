package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Telegram BotConfig
}

type BotConfig struct {
	Token string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error load .env file")
	}

	return &Config{
		Telegram: BotConfig{
			Token: os.Getenv("TOKEN"),
		},
	}
}