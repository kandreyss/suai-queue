package handlers

import "gopkg.in/telebot.v3"

func StartHandler(b *telebot.Bot) {
	b.Handle("/start", func(c telebot.Context) error {
		return c.Send("Добро пожаловать! Используйте /register для регистрации")
	})
}
