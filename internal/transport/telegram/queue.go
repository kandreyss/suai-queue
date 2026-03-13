package telegram

import (
	"fmt"
	"html"
	"strings"
	"time"

	"suai-queue/internal/domain"

	"gopkg.in/telebot.v3"
)

func (h *Handler) JoinQueue(c telebot.Context) error {
	userID := c.Sender().ID

	if !h.Repo.Exists(userID) {
		return c.Send("Сначала нужно зарегистрироваться! Введите /register")
	}

	newStudent := &domain.Student{
		TgID:          userID,
		TelegramLogin: c.Sender().Username,
		Name:          h.Repo.GetName(userID),
		TimeInQueue:   time.Now(),
	}

	position, err := h.Queue.Push(newStudent)
	if err != nil {
		return c.Send(fmt.Sprintf("Вы уже в очереди! Ваш номер: %d", position), MainMenu)
	}

	return c.Send(fmt.Sprintf("Вы успешно встали в очередь! 📝 Ваша позиция: %d", position), MainMenu)
}

func (h *Handler) LeaveQueue(c telebot.Context) error {
	userID := c.Sender().ID

	if !h.Repo.Exists(userID) {
		return c.Send("Сначала нужно зарегистрироваться! Введите /register")
	}

	err := h.Queue.Remove(userID)
	if err != nil {
		return c.Send("Вы не состоите в очереди", MainMenu)
	}

	return c.Send("Вы вышли из очереди.", MainMenu)
}

func (h *Handler) ViewQueue(c telebot.Context) error {
	users := h.Queue.GetUsers()
	if len(users) == 0 {
		return c.Send("Очередь пуста! Успей занять, пока пусто!", MainMenu)
	}

	var sb strings.Builder
	sb.WriteString("<b>Текущая очередь:</b>\n\n")

	for i, s := range users {
		name := html.EscapeString(s.Name)
		login := html.EscapeString(s.TelegramLogin)

		if c.Sender().ID == s.TgID {
			fmt.Fprintf(&sb, "<b>%d. %s @%s (Вы)</b>\n", i+1, name, login)
		} else {
			fmt.Fprintf(&sb, "%d. %s @%s\n", i+1, name, login)
		}
	}

	return c.Send(
		sb.String(),
		&telebot.SendOptions{ParseMode: telebot.ModeHTML},
		MainMenu,
	)
}
