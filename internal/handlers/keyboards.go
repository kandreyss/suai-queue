package handlers

import "gopkg.in/telebot.v3"

var (
	MainMenu = &telebot.ReplyMarkup{ResizeKeyboard: true}

	ButtonJoinInQueue = MainMenu.Text("Встать в очередь")
	ButtonViewQueue = MainMenu.Text("Посмотреть очередь")
	ButtonLeave = MainMenu.Text("Выйти из очереди")
)

var (
	SettingsMenu = &telebot.ReplyMarkup{ResizeKeyboard: true}

	ButtonChangeName = SettingsMenu.Text("Изменить имя")
)

func init() {
	MainMenu.Reply(
		MainMenu.Row(ButtonJoinInQueue, ButtonLeave),
		MainMenu.Row(ButtonViewQueue),
	)
	
	SettingsMenu.Reply(
		SettingsMenu.Row(ButtonChangeName),
	)
}