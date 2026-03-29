package telegram

import "gopkg.in/telebot.v3"

var (
	MainMenu = &telebot.ReplyMarkup{ResizeKeyboard: true}

	ButtonJoinInQueue = MainMenu.Text("➕Встать в очередь")
	ButtonViewQueue   = MainMenu.Text("📋Посмотреть очередь")
	ButtonLeave       = MainMenu.Text("➖Выйти из очереди")
	ButtonInfo        = MainMenu.Text("ℹ️Инфо")
)

var (
	SettingsMenu = &telebot.ReplyMarkup{ResizeKeyboard: true}

	ButtonChangeName  = SettingsMenu.Text("Изменить имя")
	ButtonChangeGroup = SettingsMenu.Text("Изменить группу")
)

func init() {
	MainMenu.Reply(
		MainMenu.Row(ButtonJoinInQueue, ButtonLeave),
		MainMenu.Row(ButtonViewQueue),
		MainMenu.Row(ButtonInfo),
	)

	SettingsMenu.Reply(
		SettingsMenu.Row(ButtonChangeName, ButtonChangeGroup),
	)
}
