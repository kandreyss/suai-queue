package handlers

import (
	"suai-queue/internal/service"

	"gopkg.in/telebot.v3"
)

func RegisterHandler(db *service.StudentService, b *telebot.Bot) {
	b.Handle("/register", func(c telebot.Context) error {
		user := c.Sender()
		userLogin := user.Username
		
	})
}