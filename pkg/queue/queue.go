package queue

import (
	"suai-queue/pkg/student"
	"sync"
)

type Queue struct {
	Users []student.Student
	mutex sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{
		Users: make([]student.Student, 0, 10),
	}
}

func (q *Queue) Pop() (*student.Student, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.Users) == 0 {
		return nil, ErrQueueIsEmpty
	}
	
	firstUser := q.Users[0]
	q.Users = q.Users[1:]

	return &firstUser, nil
}

func (q *Queue) Push(user *student.Student) (int, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for i, u := range q.Users {
		if u.ID == user.ID {
			return i+1, ErrStudentInQueue
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
		if user.ID == userID {
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

func (q *Queue) GetUsers() []student.Student {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	usersCopy := make([]student.Student, len(q.Users))
	copy(usersCopy, q.Users)
	return usersCopy
}