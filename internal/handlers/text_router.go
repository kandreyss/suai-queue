package handlers

import (
	"suai-queue/internal/service"

	"gopkg.in/telebot.v3"
)

func TextRouterHandler(db *service.StudentService, b *telebot.Bot) {
	b.Handle(telebot.OnText, func(c telebot.Context) error {
		userID := c.Sender().ID

		session, ok := sessions[userID]
		if !ok {
			return nil
		}

		switch session.State {
		case "waiting_name":
			return handleRegisterName(db)(c)

		case "waiting_setting", "waiting_new_name":
			return handleSetting(db)(c)
		}

		return nil
	})
}