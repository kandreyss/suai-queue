package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"suai-queue/internal/handlers"
	"suai-queue/internal/repository/students"
	"suai-queue/pkg/queue"

	"gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

type App struct {
	bot               *telebot.Bot
	studentRepository *students.StudentRepository
	queue             *queue.Queue
	db                *gorm.DB
	ctx               context.Context
	cancel            context.CancelFunc
}

// New создает новое приложение
func New() (*App, error) {
	// Инициализация БД
	db, err := setupDatabase()
	if err != nil {
		return nil, err
	}

	// Инициализация бота
	bot, err := setupBot()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	app := &App{
		bot:               bot,
		studentRepository: students.New(db),
		queue:             queue.NewQueue(),
		db:                db,
		ctx:               ctx,
		cancel:            cancel,
	}

	return app, nil
}

// Run запускает приложение
func (a *App) Run() {
	a.startServices()
	a.registerHandlers()
	a.setupShutdown()

	log.Println("Бот запущен 🚀")
	a.bot.Start()
}

// Close закрывает приложение
func (a *App) Close() {
	a.cancel()

	// Закрытие БД
	if a.db != nil {
		sqlDB, err := a.db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}

// startServices запускает все сервисы
func (a *App) startServices() {
	setupQueueCleanup(a.ctx, a.bot, a.queue)
}

// registerHandlers регистрирует все обработчики
func (a *App) registerHandlers() {
	handlers.RegisterHandler(a.studentRepository, a.bot)
	handlers.SettingsHandler(a.studentRepository, a.bot)
	handlers.TextRouterHandler(a.studentRepository, a.queue, a.bot)
	handlers.StartHandler(a.studentRepository, a.bot)
	handlers.QueueHandlers(a.studentRepository, a.queue, a.bot)
	handlers.HelpHandler(a.bot)
}

// setupShutdown настраивает graceful shutdown
func (a *App) setupShutdown() {
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		log.Println("Завершение работы...")
		a.cancel()
		a.bot.Stop()
	}()
}
