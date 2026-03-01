package handlers

import (
	"fmt"
	"strings"
	"suai-queue/internal/service"
	"suai-queue/pkg/queue"
	"suai-queue/pkg/student"

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
			c.Send("Очередь пуста! Успей занять, пока пусто!", MainMenu)
		}

		var studentsList strings.Builder
		studentsList.WriteString("*Текущая очередь:*\n\n")
		for i, s := range q.GetUsers() {
			fmt.Fprintf(&studentsList, "%d. %s @%s\n", i+1, s.Name, s.TelegramLogin)
		}

		return c.Send(studentsList.String(), &telebot.SendOptions{ParseMode: telebot.ModeMarkdown}, MainMenu)
	}
}