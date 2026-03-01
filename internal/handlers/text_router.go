package handlers

import (
	"suai-queue/internal/handlers/sessions"
	"suai-queue/internal/service"

	"gopkg.in/telebot.v3"
)

func TextRouterHandler(db *service.StudentService, b *telebot.Bot) {
	b.Handle(telebot.OnText, func(c telebot.Context) error {
		userID := c.Sender().ID

		session, ok := sessions.Store.Get(userID)
		if !ok {
			return nil
		}

		switch session.State {
		case sessions.StateWaitingName:
			return handleRegisterName(db, c, userID, session)
			
		case sessions.StateWaitingSetting, sessions.StateWaitingNewName:
			return handleSetting(db, c, userID, session)
			
		default:
			sessions.Store.Delete(userID)
			return nil
		}
	})
}