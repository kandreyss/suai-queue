package telegram

import (
	"fmt"
	"strings"

	"gopkg.in/telebot.v3"
)

func (h *Handler) Start(c telebot.Context) error {
	userID := c.Sender().ID

	if h.Repo.Exists(userID) {
		return c.Send("С возвращением! Выберите действие ниже 👇", MainMenu)
	}

	err := c.Send(helpText, telebot.ModeMarkdownV2)
	if err != nil {
		return err
	}

	return c.Send("Добро пожаловать! Пожалуйста, используйте /register для регистрации.")
}

func (h *Handler) Help(c telebot.Context) error {
	return c.Send(helpText, &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
}

func (h *Handler) Info(c telebot.Context) error {
	userID := c.Sender().ID

	if !h.Repo.Exists(userID) {
		return c.Send("Сначала нужно зарегистрироваться! Введите /register")
	}

	student, err := h.Repo.GetByTGID(userID)
	if err != nil {
		return c.Send("Ошибка получения информации из базы данных, попробуйте еще раз!", MainMenu)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Имя: <i>%s</i>\n", student.Name)
	fmt.Fprintf(&sb, "Группа: <i>%s</i>\n", student.Group)
	fmt.Fprintf(&sb, "\nЧтобы сменить информацию в профиле воспользуйтесь /settings")

	return c.Send(
		sb.String(),
		&telebot.SendOptions{ParseMode: telebot.ModeHTML},
		MainMenu,
	)
}

var helpText = `*🛰 Помощник по очередям ГУАП*

Привет\! Этот бот создан специально для студентов **СПбГУАП**, чтобы навсегда забыть про споры у дверей кафедр и списки на листочках\.

*Предназначение:*
Бот автоматизирует очередь на сдачу лабораторных, курсовых и зачетов\. Вставай в список дистанционно и следи за своим номером в реальном времени\.

*🛠 Доступные команды:*

• /start — Запустить бота и открыть главное меню\.
• /register — Регистрация \(тебе нужно указать свое имя\)\.
• /settings — Настройки профиля\.

*💡 Как пользоваться?*
1\. Пройди регистрацию через /register\.
2\. Нажми кнопку *«➕ Встать в очередь»*\.

Пожалуйста\, нажимай *«Выйти из очереди»*, если ты уже сдал работу или передумал отвечать, чтобы не задерживать других студентов\. Иначе через 30 минут бездействия, бот автоматически удалит тебя из очереди\.`
