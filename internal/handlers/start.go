package handlers

import (
	"gopkg.in/telebot.v3"
	"suai-queue/internal/service"
)

func StartHandler(db *service.StudentService, b *telebot.Bot) {
    b.Handle("/start", func(c telebot.Context) error {
        userID := c.Sender().ID

        if db.Exists(userID) {
            return c.Send("С возвращением! Выберите действие ниже 👇", MainMenu)
        }

        err := c.Send(helpText, telebot.ModeMarkdownV2)
        if err != nil {
            return err
        }

        return c.Send("Добро пожаловать! Пожалуйста, используйте /register для регистрации.")
    })
}
