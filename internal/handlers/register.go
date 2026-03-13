package handlers

import (
	"fmt"
	"suai-queue/internal/handlers/sessions"
	"suai-queue/internal/repository/students"
	"suai-queue/pkg/student"

	"gopkg.in/telebot.v3"
)

func RegisterHandler(repo *students.StudentRepository, b *telebot.Bot) {
	b.Handle("/register", func(c telebot.Context) error {
		userID := c.Sender().ID

		if repo.Exists(userID) {
			return c.Send("Вы уже зарегистрированы! Приятного использования 😊", MainMenu)
		}

		sessions.Store.Set(userID, sessions.NewUserSession(sessions.StateWaitingName))
		return c.Send("Введите ваше имя:", &telebot.ReplyMarkup{ForceReply: true})
	})
}

func HandleRegisterName(
	repo *students.StudentRepository,
	c telebot.Context,
	userID int64,
) error {

	name := c.Text()
	if len([]rune(name)) < 2 {
		return c.Send("Имя слишком короткое. Введите корректное имя:")
	}

	username := c.Sender().Username
	if username == "" {
		username = "NoUsername"
	}

	st := &student.Student{
		TgID:          userID,
		TelegramLogin: username,
		Name:          name,
	}

	if err := repo.Create(st); err != nil {
		return c.Send("Ошибка при сохранении данных в базу. Попробуйте позже.")
	}

	sessions.Store.Delete(userID)
	return c.Send(
		fmt.Sprintf("Регистрация завершена! Добро пожаловать, %s", st.Name),
		MainMenu,
	)
}
