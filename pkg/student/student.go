package student

import "time"

type Student struct {
	ID            int64
	TelegramLogin string
	Name          string
	TimeToQueue   time.Time
}

func NewStudent(id int64, login string, name string) *Student {
	return &Student{
		ID:            id,
		TelegramLogin: login,
		Name:          name,
	}
}
