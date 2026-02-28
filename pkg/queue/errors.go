package queue

import "errors"

var (
	ErrQueueIsEmpty = errors.New("queue is empty. nothing to pop()")
)
