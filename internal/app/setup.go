package app

import (
	"context"
	"time"

	"suai-queue/internal/config"
	"suai-queue/internal/database"
	"suai-queue/internal/service/queue"

	"gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

func setupBot() (*telebot.Bot, error) {
	cfg := config.MustLoad()

	bot, err := telebot.NewBot(telebot.Settings{
		Token:  cfg.Telegram.Token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, err
	}

	if err := setBotCommands(bot); err != nil {
		return nil, err
	}

	return bot, nil
}

func setBotCommands(b *telebot.Bot) error {
	return b.SetCommands([]telebot.Command{
		{Text: "start", Description: "Запуск и главное меню"},
		{Text: "register", Description: "Регистрация"},
		{Text: "settings", Description: "Настройки"},
		{Text: "help", Description: "Список команд и помощь"},
	})
}

func setupDatabase() (*gorm.DB, error) {
	cfg := config.MustLoad()
	return database.InitDB(cfg)
}

func setupQueueCleanup(ctx context.Context, bot *telebot.Bot, q *queue.Queue) {
	queue.StartQueueCleanup(ctx, bot, q, 10*time.Minute, 25*time.Minute)
}
