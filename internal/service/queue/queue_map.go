package queue

import (
	"strings"
	"sync"
)

type QueueMap struct {
	mutex  sync.Mutex
	Queues map[string]*Queue
}

func NewQueueMap() *QueueMap {
	return &QueueMap{
		Queues: make(map[string]*Queue),
	}
}

func (qm *QueueMap) GetQueue(groupNumber string) (*Queue, bool) {
	qm.mutex.Lock()
	defer qm.mutex.Unlock()

	queue, exists := qm.Queues[groupNumber]
	return queue, exists
}

func (qm *QueueMap) RequireQueue(groupNumber string) (*Queue, error) {
	normalizedGroup := strings.TrimSpace(groupNumber)
	if normalizedGroup == "" {
		return nil, ErrInvalidGroup
	}

	qm.mutex.Lock()
	defer qm.mutex.Unlock()

	queue, exists := qm.Queues[normalizedGroup]
	if !exists {
		return nil, ErrQueueNotFound
	}

	return queue, nil
}

func (qm *QueueMap) EnsureQueue(groupNumber string) (*Queue, error) {
	normalizedGroup := strings.TrimSpace(groupNumber)
	if normalizedGroup == "" {
		return nil, ErrInvalidGroup
	}

	qm.mutex.Lock()
	defer qm.mutex.Unlock()

	if queue, exists := qm.Queues[normalizedGroup]; exists {
		return queue, nil
	}

	queue := NewQueue()
	qm.Queues[normalizedGroup] = queue
	return queue, nil
}

func (qm *QueueMap) AddQueue(groupNumber string, queue *Queue) {
	normalizedGroup := strings.TrimSpace(groupNumber)
	if normalizedGroup == "" || queue == nil {
		return
	}

	qm.mutex.Lock()
	defer qm.mutex.Unlock()

	qm.Queues[normalizedGroup] = queue
}

func (qm *QueueMap) RemoveQueue(groupNumber string) {
	normalizedGroup := strings.TrimSpace(groupNumber)
	if normalizedGroup == "" {
		return
	}

	qm.mutex.Lock()
	defer qm.mutex.Unlock()

	delete(qm.Queues, normalizedGroup)
}
