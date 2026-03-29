package telegram

import (
	"errors"
	"suai-queue/internal/domain"
	"suai-queue/internal/repository"
	"suai-queue/internal/service/queue"
	"suai-queue/internal/session"
	"time"

	"gopkg.in/telebot.v3"
)

var acceptableSettings = []string{
	"Изменить имя",
	"Изменить группу",
}

func isAcceptableConfiguration(configuration string) bool {
	for _, conf := range acceptableSettings {
		if configuration == conf {
			return true
		}
	}
	return false
}

func (h *Handler) Settings(c telebot.Context) error {
	userID := c.Sender().ID

	if !h.Repo.Exists(userID) {
		return c.Send("Для начала, используйте /register для регистрации")
	}

	h.Session.Set(userID, session.NewUserSession(session.StateWaitingSetting))
	return c.Send("Выберите необходимую настройку ниже👇", SettingsMenu)
}

func (h *Handler) handleSettingUpdate(c telebot.Context, userID int64, sess *session.UserSession) error {
	switch sess.State {
	case session.StateWaitingSetting:
		setting := c.Text()

		if !isAcceptableConfiguration(setting) {
			return c.Send("Неизвестная настройка! Воспользуйтесь настройками ниже👇", SettingsMenu)
		}

		switch setting {
		case "Изменить имя":
			sess.State = session.StateWaitingNewName
			h.Session.Set(userID, sess)
			return c.Send("Введите новое имя:", &telebot.ReplyMarkup{RemoveKeyboard: true})

		case "Изменить группу":
			sess.State = session.StateWaitingNewGroup
			h.Session.Set(userID, sess)
			return c.Send("Введите новый номер группы:", &telebot.ReplyMarkup{RemoveKeyboard: true})

		default:
			return c.Send("Эта настройка пока не реализована.", SettingsMenu)
		}

	case session.StateWaitingNewName:
		newName, err := readAndValidateName(c)
		if err != nil {
			return c.Send("Имя слишком короткое. Введите корректное имя:")
		}

		if err := h.Repo.UpdateName(userID, newName); err != nil {
			h.Session.Delete(userID)
			return c.Send("Не удалось обновить имя. Попробуйте позже.", MainMenu)
		}

		student, err := h.Repo.GetByTGID(userID)
		if err == nil {
			if groupQueue, qErr := h.Queues.RequireQueue(student.Group); qErr == nil {
				groupQueue.UpdateQueueUser(userID, func(s *domain.Student) {
					s.Name = newName
				})
			} else if !errors.Is(qErr, queue.ErrQueueNotFound) {
				h.Session.Delete(userID)
				return c.Send("Не удалось синхронизировать данные в очереди. Попробуйте позже.", MainMenu)
			}
		}

		h.Session.Delete(userID)
		return c.Send("Имя успешно обновлено ✅", MainMenu)

	case session.StateWaitingNewGroup:
		newGroup, err := readAndValidateGroup(c)
		if err != nil {
			return c.Send("Номер группы должен состоять минимум из 4-х символов.")
		}

		student, err := h.Repo.GetByTGID(userID)
		if err != nil {
			h.Session.Delete(userID)
			return c.Send("Не удалось получить данные профиля. Попробуйте позже.", MainMenu)
		}

		oldGroup := student.Group

		if err := h.Repo.UpdateGroup(userID, newGroup); err != nil {
			h.Session.Delete(userID)
			if errors.Is(err, repository.ErrGroupNotFound) {
				return c.Send("Группа не найдена! Проверьте номер и попробуйте ещё раз", MainMenu)
			}
			return c.Send("Не удалось обновить группу. Попробуйте позже.", MainMenu)
		}

		oldQueue, qErr := h.Queues.RequireQueue(oldGroup)
		switch {
		case qErr == nil:
			if rmErr := oldQueue.Remove(userID); rmErr == nil {
				updatedStudent := &domain.Student{
					TgID:          student.TgID,
					TelegramLogin: student.TelegramLogin,
					Name:          student.Name,
					Group:         newGroup,
					TimeInQueue:   time.Now(),
				}

				newQueue, ensureErr := h.Queues.EnsureQueue(newGroup)
				if ensureErr != nil {
					h.Session.Delete(userID)
					return c.Send("Группа обновлена, но не удалось перенести в очередь новой группы.", MainMenu)
				}

				if _, pushErr := newQueue.Push(updatedStudent); pushErr != nil && !errors.Is(pushErr, queue.ErrStudentInQueue) {
					h.Session.Delete(userID)
					return c.Send("Группа обновлена, но не удалось перенести в очередь новой группы.", MainMenu)
				}
			}
		case errors.Is(qErr, queue.ErrQueueNotFound):
		default:
			h.Session.Delete(userID)
			return c.Send("Группа обновлена, но не удалось синхронизировать очередь.", MainMenu)
		}

		h.Session.Delete(userID)
		return c.Send("Группа успешно обновлена ✅", MainMenu)
	}

	return nil
}
