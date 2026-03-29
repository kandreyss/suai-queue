package telegram

import (
	"errors"
	"fmt"
	"html"
	"strings"
	"time"

	"suai-queue/internal/domain"
	"suai-queue/internal/service/queue"

	"gopkg.in/telebot.v3"
)

func (h *Handler) JoinQueue(c telebot.Context) error {
	userID := c.Sender().ID

	student, err := h.Repo.GetByTGID(userID)
	if err != nil {
		return c.Send("Сначала нужно зарегистрироваться! Введите /register")
	}

	newStudent := &domain.Student{
		TgID:          userID,
		TelegramLogin: c.Sender().Username,
		Name:          student.Name,
		Group:         student.Group,
		TimeInQueue:   time.Now(),
	}

	groupQueue, err := h.Queues.EnsureQueue(student.Group)
	if err != nil {
		return c.Send("Не удалось определить очередь вашей группы.", MainMenu)
	}

	position, err := groupQueue.Push(newStudent)
	if err != nil {
		if errors.Is(err, queue.ErrStudentInQueue) {
			return c.Send(fmt.Sprintf("Вы уже в очереди! Ваш номер: %d", position), MainMenu)
		}
		return c.Send("Не удалось добавить вас в очередь. Попробуйте позже.", MainMenu)
	}

	return c.Send(fmt.Sprintf("Вы успешно встали в очередь! 📝 Ваша позиция: %d", position), MainMenu)
}

func (h *Handler) LeaveQueue(c telebot.Context) error {
	userID := c.Sender().ID

	student, err := h.Repo.GetByTGID(userID)
	if err != nil {
		return c.Send("Сначала нужно зарегистрироваться! Введите /register")
	}

	groupQueue, err := h.Queues.RequireQueue(student.Group)
	if err != nil {
		if errors.Is(err, queue.ErrQueueNotFound) {
			return c.Send("Вы не состоите в очереди", MainMenu)
		}
		return c.Send("Не удалось определить очередь вашей группы.", MainMenu)
	}

	err = groupQueue.Remove(userID)
	if err != nil {
		return c.Send("Вы не состоите в очереди", MainMenu)
	}

	return c.Send("Вы вышли из очереди.", MainMenu)
}

func (h *Handler) ViewQueue(c telebot.Context) error {
	student, err := h.Repo.GetByTGID(c.Sender().ID)
	if err != nil {
		return c.Send("Сначала нужно зарегистрироваться! Введите /register")
	}

	groupQueue, err := h.Queues.RequireQueue(student.Group)
	if err != nil {
		if errors.Is(err, queue.ErrQueueNotFound) {
			return c.Send("Очередь пуста! Успей занять, пока пусто!", MainMenu)
		}
		return c.Send("Не удалось определить очередь вашей группы.", MainMenu)
	}

	users := groupQueue.GetUsers()
	if len(users) == 0 {
		return c.Send("Очередь пуста! Успей занять, пока пусто!", MainMenu)
	}

	var sb strings.Builder
	sb.WriteString("<i>Группа: " + student.Group + "</i>\n<b>Текущая очередь: </b>\n\n")

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
