package sessions

type UserSession struct {
	State    string
	TempData map[string]string
}

func NewUserSession(state string) *UserSession {
	return &UserSession{
		State:    state,
		TempData: make(map[string]string),
	}
}

var Store = NewSessionStore()
