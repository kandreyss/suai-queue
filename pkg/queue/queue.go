package queue

import (
	"suai-queue/pkg/student"
	"sync"
)

type Queue struct {
	Users []student.Student
	mutex sync.Mutex
}

func (q *Queue) New() *Queue {
	return &Queue{
		Users: make([]student.Student, 0, 10),
	}
}

func (q *Queue) Pop() (*student.Student, error) {
	if(len(q.Users) == 0) {
		return nil, ErrQueueIsEmpty
	}
	q.mutex.Lock()
	firstUser := q.Users[0]
	q.Users = q.Users[1:]
	q.mutex.Unlock()

	return &firstUser, nil
} 

func (q *Queue) Push(user student.Student) {
	q.Users = append(q.Users, user)
}

