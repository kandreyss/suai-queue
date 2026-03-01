package handlers

import (
	"strings"

	"suai-queue/internal/handlers/sessions"
	"suai-queue/internal/service"
	"suai-queue/pkg/queue"

	"gopkg.in/telebot.v3"
)

func TextRouterHandler(db *service.StudentService, q *queue.Queue, b *telebot.Bot) {
	b.Handle(telebot.OnText, func(c telebot.Context) error {
		userID := c.Sender().ID
		_ = strings.TrimSpace(c.Text())

		if session, ok := sessions.Store.Get(userID); ok {
			switch session.State {
			case sessions.StateWaitingName:
				return handleRegisterName(db, c, userID, session)

			case sessions.StateWaitingSetting, sessions.StateWaitingNewName:
				return handleSetting(db, q, c, userID, session)

			default:

				sessions.Store.Delete(userID)
			}
		}

		if !db.Exists(userID) {
			return c.Send("Для использования бота, необходимо зарегистрироваться с помощью /register")
		}

		return c.Send("Не понял сообщение. Используй кнопки меню 👇", MainMenu)
	})
}
