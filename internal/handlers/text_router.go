package handlers

import (
	"suai-queue/internal/handlers/sessions"
	"suai-queue/internal/service"

	telebot "gopkg.in/telebot.v3"
)

func TextRouterHandler(db *service.StudentService, b *telebot.Bot) {
	b.Handle(telebot.OnText, func(c telebot.Context) error {
		userID := c.Sender().ID

		session, ok := (*sessions.SessionStore).Get(userID)
		if !ok {
			return nil
		}

		switch session.State {
		case session.StateWaitingName:
			return handleRegisterNameStep(db, c, userID, session)
		default:
			sessionsStore.Delete(userID)
			return nil
		}
	})
}