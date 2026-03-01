package handlers

import (
	"fmt"

	"suai-queue/internal/service"
	"suai-queue/pkg/student"
	
	"gopkg.in/telebot.v3"
)

func RegisterHandler(db *service.StudentService, b *telebot.Bot) {
	b.Handle("/register", handleRegisterCommand(db))
}

func handleRegisterCommand(db *service.StudentService) func(telebot.Context) error {
	return func(c telebot.Context) error {
		userID := c.Sender().ID

		if db.Exists(userID) {
			return c.Send("Вы уже зарегистрированы! Приятного использования 😊", MainMenu)
		}

		sessionsStore.Set(userID, NewUserSession(StateWaitingName))
		return c.Send("Введите ваше имя:", &telebot.ReplyMarkup{ForceReply: true})
	}
}

func handleRegisterNameStep(db *service.StudentService, c telebot.Context, userID int64, session *UserSession) error {
	name, err := readAndValidateName(c)
	if err != nil {
		if err == ErrNameTooShort {
			return c.Send("Имя слишком короткое. Введите корректное имя:")
		}
		return c.Send("Некорректное имя. Попробуйте ещё раз:")
	}

	st, err := registerStudent(db, c, userID, name)
	if err != nil {
		if err == service.ErrStudentInDb {
			endSession(userID)
			return c.Send("Вы уже зарегистрированы!", MainMenu)
		}
		return c.Send("Ошибка при сохранении данных. Попробуйте позже.")
	}

	endSession(userID)
	return c.Send(
		fmt.Sprintf("Регистрация завершена! Добро пожаловать, %s", st.Name),
		MainMenu,
	)
}

func registerStudent(db *service.StudentService, c telebot.Context, userID int64, name string) (*student.Student, error) {
	username := c.Sender().Username
	if username == "" {
		username = "NoUsername"
	}

	st := student.NewStudent(userID, username, name)
	if err := db.Insert(st); err != nil {
		return nil, err
	}
	return st, nil
}