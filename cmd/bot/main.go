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

	studentService := service.NewStudentService()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				students := studentService.GetStudents()

				log.Println("==== База студентов ====")
				for _, s := range students {
					log.Printf("ID: %d | Name: %s\n", s.ID, s.Name)
				}
				log.Println("========================")

			case <-ctx.Done():
				log.Println("Фоновая горутина завершена")
				return
			}
		}
	}()

	handlers.RegisterHandler(studentService, bot)
	handlers.SettingsHandler(studentService, bot)
	handlers.TextRouterHandler(studentService, bot)
	handlers.StartHandler(studentService, bot)

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
