package handlers

import (
	"fmt"

	"suai-queue/internal/service"
	"suai-queue/pkg/student"

	"gopkg.in/telebot.v3"
)

func RegisterHandler(db *service.StudentService, b *telebot.Bot) {
	b.Handle("/register", func(c telebot.Context) error {
		userID := c.Sender().ID

		if db.Exists(userID) {
			return c.Send("Вы уже зарегистрированы! Приятного использования 😊", MainMenu)
		}

		sessionsStore.Set(userID, NewUserSession(StateWaitingName))
		return c.Send("Введите ваше имя:", &telebot.ReplyMarkup{ForceReply: true})
	})
}

func handleRegisterName(db *service.StudentService, c telebot.Context, userID int64, session *UserSession) error {
	name := c.Text()
	if len([]rune(name)) < 2 {
		return c.Send("Имя слишком короткое. Введите корректное имя:")
	}

	username := c.Sender().Username
	if username == "" {
		username = "NoUsername"
	}

	st := student.NewStudent(userID, username, name)

	if err := db.Insert(st); err != nil {
		if err == service.ErrStudentInDb {
			sessionsStore.Delete(userID)
			return c.Send("Вы уже зарегистрированы!", MainMenu)
		}
		return c.Send("Ошибка при сохранении данных. Попробуйте позже.")
	}

	sessionsStore.Delete(userID)
	return c.Send(
		fmt.Sprintf("Регистрация завершена! Добро пожаловать, %s", st.Name),
		MainMenu,
	)
}