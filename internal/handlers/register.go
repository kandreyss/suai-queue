package handlers

import (
	"fmt"

	"suai-queue/internal/service"
	"suai-queue/pkg/student"

	"gopkg.in/telebot.v3"
)

var sessions = make(map[int64]*UserSession)

func RegisterHandler(db *service.StudentService, b *telebot.Bot) {

	b.Handle("/register", func(c telebot.Context) error {
		userID := c.Sender().ID

		if db.Exists(userID) {
			return c.Send("Вы уже зарегистрированы! Приятного использования 😊", MainMenu)
		}

		sessions[userID] = NewUserSession("waiting_name")

		return c.Send("Введите ваше имя:", &telebot.ReplyMarkup{
			ForceReply: true,
		})
	})

	b.Handle(telebot.OnText, func(c telebot.Context) error {
		userID := c.Sender().ID
		session, ok := sessions[userID]

		if !ok {
			return nil
		}

		switch session.State {

		case "waiting_name":
			name := c.Text()

			if len(name) < 2 {
				return c.Send("Имя слишком короткое. Введите корректное имя:")
			}

			st := student.NewStudent(
				userID,
				c.Sender().Username,
				name,
			)

			if err := db.Insert(st); err != nil {
				if err == service.ErrStudentInDb {
					delete(sessions, userID)
					return c.Send("Вы уже зарегистрированы!")
				}
				return c.Send("Ошибка при сохранении данных. Попробуйте позже.")
			}

			delete(sessions, userID)

			return c.Send(
				fmt.Sprintf("Регистрация завершена! Добро пожаловать, %s", st.Name),
				MainMenu,
			)
		}

		return nil
	})
}