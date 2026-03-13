package student

import (
	"time"

	"gorm.io/gorm"
)

type Student struct {
	gorm.Model

	Group         uint8  `gorm:"column:group_num"`
	TgID          int64  `gorm:"column:tg_id;uniqueIndex"`
	TelegramLogin string `gorm:"column:tg_login"`
	Name          string `gorm:"column:name"`

	TimeInQueue time.Time
}

func NewStudent(id int64, group uint8, login string, name string) *Student {
	return &Student{
		Group:         group,
		TgID:          id,
		TelegramLogin: login,
		Name:          name,
	}
}
