package handlers

import (
	"suai-queue/internal/handlers/sessions"
	"suai-queue/internal/service"
	"suai-queue/pkg/queue"

	"gopkg.in/telebot.v3"
)

var acceptableSettings []string = []string{
	"Изменить имя",
}

func isAcceptableConfiguration(configuration string) bool {
	for _, conf := range acceptableSettings {
		if configuration == conf {
			return true
		}
	}
	return false
}

func SettingsHandler(db *service.StudentService, b *telebot.Bot) {
	b.Handle("/settings", func(c telebot.Context) error {
		userID := c.Sender().ID

		if !db.Exists(userID) {
			return c.Send("Для начала, используйте /register для регистрации")
		}

		sessions.Store.Set(userID, sessions.NewUserSession(sessions.StateWaitingSetting))
		return c.Send("Выберите необходимую настройку ниже👇", SettingsMenu)
	})
}

func handleSetting(db *service.StudentService, q *queue.Queue, c telebot.Context, userID int64, session *sessions.UserSession) error {
	switch session.State {
	case sessions.StateWaitingSetting:
		setting := c.Text()

		if !isAcceptableConfiguration(setting) {
			return c.Send("Неизвестная настройка! Воспользуйтесь настройками ниже👇", SettingsMenu)
		}

		switch setting {
		case "Изменить имя":
			session.State = sessions.StateWaitingNewName
			sessions.Store.Set(userID, session)
			return c.Send("Введите новое имя:")

		default:
			return c.Send("Эта настройка пока не реализована.", SettingsMenu)
		}

	case sessions.StateWaitingNewName:
		newName, err := readAndValidateName(c)
		if err != nil {
			if err == ErrNameTooShort {
				return c.Send("Имя слишком короткое. Введите корректное имя:")
			}
			return c.Send("Некорректное имя. Попробуйте ещё раз:")
		}

		if err := db.UpdateName(userID, newName); err != nil {
			sessions.Store.Delete(userID)
			return c.Send("Не удалось обновить имя. Используйте команду /settings заново.", MainMenu)
		}

		sessions.Store.Delete(userID)
		for i := range q.Users {
			if q.Users[i].ID == userID {
				q.Users[i].Name = newName
			}
		}
		return c.Send("Имя успешно обновлено ✅", MainMenu)
	}

	return nil
}
