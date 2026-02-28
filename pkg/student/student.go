package student

type Student struct {
	ID            int64
	TelegramLogin string
	Name          string
}

func NewStudent(id int64, login string, name string) *Student {
	return &Student{
		ID:            id,
		TelegramLogin: login,
		Name:          name,
	}
}
