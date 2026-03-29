package telegram

import (
	"errors"
	"strings"

	"gopkg.in/telebot.v3"
	"suai-queue/internal/session"
)

var (
	ErrNameTooShort  = errors.New("name too short")
	ErrGroupTooShort = errors.New("group number too short")
)

func readAndValidateName(c telebot.Context) (string, error) {
	name := strings.TrimSpace(c.Text())
	if len([]rune(name)) < 2 {
		return "", ErrNameTooShort
	}
	return name, nil
}

func readAndValidateGroup(c telebot.Context) (string, error) {
	group := strings.TrimSpace(c.Text())
	if len([]rune(group)) < 4 {
		return "", ErrGroupTooShort
	}
	return group, nil
}

func (h *Handler) OnText(c telebot.Context) error {
	userID := c.Sender().ID

	if sess, ok := h.Session.Get(userID); ok {
		switch sess.State {
		case session.StateWaitingName:
			return h.handleRegisterName(c, userID, sess)

		case session.StateWaitingGroup:
			return h.handleRegisterGroup(c, userID, sess)

		case session.StateWaitingSetting, session.StateWaitingNewName, session.StateWaitingNewGroup:
			return h.handleSettingUpdate(c, userID, sess)

		default:
			h.Session.Delete(userID)
		}
	}

	if !h.Repo.Exists(userID) {
		return c.Send("Для использования бота, необходимо зарегистрироваться с помощью /register")
	}

	return c.Send("Не понял сообщение. Используй кнопки меню 👇", MainMenu)
}
