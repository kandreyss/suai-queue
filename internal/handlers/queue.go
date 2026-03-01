package handlers

import (
	"fmt"
	"strings"
	"suai-queue/internal/service"
	"suai-queue/pkg/queue"
	"suai-queue/pkg/student"
	"time"
	"html"

	"gopkg.in/telebot.v3"
)

func QueueHandlers(db *service.StudentService, q *queue.Queue, b *telebot.Bot) {
	b.Handle(&ButtonJoinInQueue, handleJoinQueue(db, q))
	b.Handle(&ButtonLeave, handleLeaveQueue(db, q))
	b.Handle(&ButtonViewQueue, handleViewQueue(q))
}

func handleJoinQueue(db *service.StudentService, q *queue.Queue) func(telebot.Context) error {
	return func(c telebot.Context) error {
		userID := c.Sender().ID

		if !db.Exists(userID) {
			return c.Send("Сначала нужно зарегистрироваться! Введите /register")
		}

		newStudent := student.NewStudent(userID, c.Sender().Username, db.GetName(userID))
		newStudent.TimeInQueue = time.Now()
		position, err := q.Push(newStudent)
		if err != nil {
			return c.Send(fmt.Sprintf("Вы уже в очереди! Ваш номер: %d", position), MainMenu)
		}

		return c.Send(fmt.Sprintf("Вы успешно встали в очередь! 📝 Ваша позиция: %d", position), MainMenu)
	}
}

func handleLeaveQueue(db *service.StudentService, q *queue.Queue) func(telebot.Context) error {
	return func(c telebot.Context) error {
		userID := c.Sender().ID

		if !db.Exists(userID) {
			return c.Send("Сначала нужно зарегистрироваться! Введите /register")
		}

		err := q.Remove(userID)
		if err != nil {
			return c.Send("Вы не состоите в очереди", MainMenu)
		}

		return c.Send("Вы вышли из очереди.", MainMenu)
	}
}

func handleViewQueue(q *queue.Queue) func(telebot.Context) error {
	return func(c telebot.Context) error {
		users := q.GetUsers()
		if len(users) == 0 {
			return c.Send("Очередь пуста! Успей занять, пока пусто!", MainMenu)
		}

		var sb strings.Builder
		sb.WriteString("<b>Текущая очередь:</b>\n\n")

		for i, s := range users {
			name := html.EscapeString(s.Name)
			login := html.EscapeString(s.TelegramLogin)

			if c.Sender().ID == s.ID {
				fmt.Fprintf(&sb, "<b>%d. %s @%s (Вы)</b>\n", i+1, name, login)
			} else {
				fmt.Fprintf(&sb, "%d. %s @%s\n", i+1, name, login)
			}
		}

		return c.Send(sb.String(), &telebot.SendOptions{ParseMode: telebot.ModeHTML}, MainMenu)
	}
}
