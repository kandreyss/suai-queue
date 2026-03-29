package telegram

import (
	"suai-queue/internal/repository"
	"suai-queue/internal/service/queue"
	"suai-queue/internal/session"

	"gopkg.in/telebot.v3"
)

type Handler struct {
	Bot     *telebot.Bot
	Repo    *repository.StudentRepository
	Queues   *queue.QueueMap
	Session *session.SessionStore
}

func NewHandler(bot *telebot.Bot, repo *repository.StudentRepository, q *queue.QueueMap, sess *session.SessionStore) *Handler {
	return &Handler{
		Bot:     bot,
		Repo:    repo,
		Queues:   q,
		Session: sess,
	}
}

func (h *Handler) Init() {
	// Commands
	h.Bot.Handle("/start", h.Start)
	h.Bot.Handle("/help", h.Help)
	h.Bot.Handle("/register", h.Register)
	h.Bot.Handle("/settings", h.Settings)
	h.Bot.Handle(&ButtonInfo, h.Info)

	// Queue
	h.Bot.Handle(&ButtonJoinInQueue, h.JoinQueue)
	h.Bot.Handle(&ButtonLeave, h.LeaveQueue)
	h.Bot.Handle(&ButtonViewQueue, h.ViewQueue)

	// Router
	h.Bot.Handle(telebot.OnText, h.OnText)
}
