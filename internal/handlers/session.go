package handlers

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

var sessionsStore = NewSessionStore()

func endSession(userID int64) {
	sessionsStore.Delete(userID)
}