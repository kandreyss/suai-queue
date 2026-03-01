package handlers

import (
	"suai-queue/internal/service"
	"suai-queue/internal/handlers/sessions"

	"gopkg.in/telebot.v3"
)

func SettingsHandler(db *service.StudentService, b *telebot.Bot) {
	b.Handle("/settings", func(c telebot.Context) error {
		userID := c.Sender().ID

		if !db.Exists(userID) {
			return c.Send("Для начала, используйте /register для регистрации")
		}

		(*sessions.SessionStore).Set(userID, sessions.NewUserSession(sessions.StateWaitingSetting))
		return c.Send("Выберите необходимую настройку ниже👇", SettingsMenu)
	})
}