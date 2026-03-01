package handlers

import (
	"errors"
	"strings"

	"gopkg.in/telebot.v3"
)

var ErrNameTooShort = errors.New("name too short")

func readAndValidateName(c telebot.Context) (string, error) {
	name := strings.TrimSpace(c.Text())
	if len([]rune(name)) < 2 {
		return "", ErrNameTooShort
	}
	return name, nil
}