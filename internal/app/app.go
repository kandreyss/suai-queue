package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"suai-queue/internal/repository"
	"suai-queue/internal/service/queue"
	"suai-queue/internal/session"
	"suai-queue/internal/transport/telegram"

	"gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

type App struct {
	bot               *telebot.Bot
	studentRepository *repository.StudentRepository
	queues            *queue.QueueMap
	sessionStore      *session.SessionStore
	db                *gorm.DB
	ctx               context.Context
	cancel            context.CancelFunc
}

func New() (*App, error) {
	db, err := setupDatabase()
	if err != nil {
		return nil, err
	}

	bot, err := setupBot()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	app := &App{
		bot:               bot,
		studentRepository: repository.New(db),
		queues:            queue.NewQueueMap(),
		sessionStore:      session.NewSessionStore(),
		db:                db,
		ctx:               ctx,
		cancel:            cancel,
	}

	return app, nil
}

func (a *App) Run() {
	a.startServices()
	a.registerHandlers()
	a.setupShutdown()

	log.Println("Бот запущен 🚀")
	a.bot.Start()
}

func (a *App) Close() {
	a.cancel()

	if a.db != nil {
		sqlDB, err := a.db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}

func (a *App) startServices() {
	setupQueueCleanup(a.ctx, a.bot, a.queues)
}

func (a *App) registerHandlers() {
	h := telegram.NewHandler(a.bot, a.studentRepository, a.queues, a.sessionStore)
	h.Init()
}

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
