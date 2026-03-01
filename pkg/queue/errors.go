package queue

import "errors"

var (
	ErrQueueIsEmpty      = errors.New("queue is empty. nothing to pop()")
	ErrStudentNotInQueue = errors.New("student is not in the queue")
	ErrStudentInQueue    = errors.New("student already in the queue")
)
