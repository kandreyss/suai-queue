package telegram

import (
	"errors"
	"fmt"

	"suai-queue/internal/domain"
	"suai-queue/internal/repository"
	"suai-queue/internal/session"

	"gopkg.in/telebot.v3"
)

func (h *Handler) Register(c telebot.Context) error {
	userID := c.Sender().ID

	if h.Repo.Exists(userID) {
		return c.Send("Вы уже зарегистрированы! Приятного использования 😊", MainMenu)
	}

	h.Session.Set(userID, session.NewUserSession(session.StateWaitingName))
	return c.Send("Введите ваше имя:", &telebot.ReplyMarkup{ForceReply: true})
}

func (h *Handler) handleRegisterName(c telebot.Context, userID int64, sess *session.UserSession) error {
	name, err := readAndValidateName(c)
	if err != nil {
		return c.Send("Имя слишком короткое. Введите корректное имя:")
	}

	sess.TempData["name"] = name
	sess.State = session.StateWaitingGroup
	h.Session.Set(userID, sess)

	return c.Send("Отлично! Теперь введите номер вашей группы:")
}

func (h *Handler) handleRegisterGroup(c telebot.Context, userID int64, sess *session.UserSession) error {
	groupNum, err := readAndValidateGroup(c)
	if err != nil {
		return c.Send("Номер группы слишком короткий. Попробуйте еще раз:")
	}

	name, ok := sess.TempData["name"]
	if !ok {
		sess.State = session.StateWaitingName
		h.Session.Set(userID, sess)
		return c.Send("Произошла ошибка. Давайте начнем сначала. Введите ваше имя:")
	}

	username := c.Sender().Username
	if username == "" {
		username = "NoUsername"
	}

	st := &domain.Student{
		TgID:          userID,
		TelegramLogin: username,
		Name:          name,
		Group:         groupNum,
	}

	if err := h.Repo.Create(st); err != nil {
		if errors.Is(err, repository.ErrGroupNotFound) {
			return c.Send("Группа не найдена! Проверьте номер и попробуйте ещё раз:")
		}
		return c.Send("Ошибка при сохранении данных в базу. Попробуйте позже.")
	}

	h.Session.Delete(userID)

	return c.Send(
		fmt.Sprintf("Регистрация завершена!\nИмя: %s\nГруппа: %s", st.Name, st.Group),
		MainMenu,
	)
}
