package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"suai-queue/internal/config"
	"suai-queue/internal/handlers"
	"suai-queue/internal/service"
	"suai-queue/pkg/queue"

	"gopkg.in/telebot.v3"
)

func main() {
	cfg := config.LoadConfig()

	pref := telebot.Settings{
		Token:  cfg.Telegram.Token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal("Ошибка создания бота:", err)
	}

	if err := setBotCommands(bot); err != nil {
		log.Fatal("Не удалось установить команды:", err)
	}

	studentService := service.NewStudentService()
	q := queue.NewQueue()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queue.StartQueueCleanup(ctx, bot, q, 10*time.Minute, 25*time.Minute)

	handlers.RegisterHandler(studentService, bot)
	handlers.SettingsHandler(studentService, bot)
	handlers.TextRouterHandler(studentService, bot)
	handlers.StartHandler(studentService, bot)
	handlers.QueueHandlers(studentService, q, bot)
	handlers.HelpHandler(bot)

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop

		log.Println("Завершение работы...")
		cancel()
		bot.Stop()
	}()

	log.Println("Бот запущен 🚀")
	bot.Start()
}

func setBotCommands(b *telebot.Bot) error {
	return b.SetCommands([]telebot.Command{
		{Text: "start", Description: "Запуск и главное меню"},
		{Text: "register", Description: "Регистрация"},
		{Text: "settings", Description: "Настройки"},
		{Text: "help", Description: "Список команд и помощь"},
	})
}