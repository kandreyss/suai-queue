package domain

import "time"

type Student struct {
	TgID          int64  `gorm:"primaryKey;autoIncrement:false"`
	TelegramLogin string `gorm:"column:tg_login"`
	Name          string `gorm:"column:name"`
	Group         string `gorm:"column:group_number;type:varchar(7)"`
	TimeInQueue   time.Time `gorm:"-"`
}
