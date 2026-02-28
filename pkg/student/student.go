package student

import "time"

type Student struct {
	TelegramLogin string
	Name string
	LoginTime time.Time
}

func (s *Student) New(login string, name string) *Student {
	return &Student{
		TelegramLogin: login,
		Name: name,
		LoginTime: time.Now(),
	}
}