package queue

import (
	"suai-queue/internal/domain"
	"sync"
)

type Queue struct {
	Users []domain.Student
	mutex sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{
		Users: make([]domain.Student, 0, 10),
	}
}

func (q *Queue) Pop() (*domain.Student, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.Users) == 0 {
		return nil, ErrQueueIsEmpty
	}

	firstUser := q.Users[0]
	q.Users = q.Users[1:]

	return &firstUser, nil
}

func (q *Queue) Push(user *domain.Student) (int, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for i, u := range q.Users {
		if u.TgID == user.TgID {
			return i + 1, ErrStudentInQueue
		}
	}

	q.Users = append(q.Users, *user)

	return len(q.Users), nil
}

func (q *Queue) Remove(userID int64) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	indexToRemove := -1
	for i, user := range q.Users {
		if user.TgID == userID {
			indexToRemove = i
			break
		}
	}

	if indexToRemove == -1 {
		return ErrStudentNotInQueue
	}

	q.Users = append(q.Users[:indexToRemove], q.Users[indexToRemove+1:]...)

	return nil
}

func (q *Queue) GetUsers() []domain.Student {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	usersCopy := make([]domain.Student, len(q.Users))
	copy(usersCopy, q.Users)
	return usersCopy
}

func (q *Queue) UpdateQueueUser(userID int64, updateFn func(*domain.Student)) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for i := range q.Users {
		if q.Users[i].TgID == userID {
			updateFn(&q.Users[i])
			break
		}
	}
}
