package queue

import (
	"context"
	"fmt"
	"log"
	"time"

	"suai-queue/pkg/student"

	"gopkg.in/telebot.v3"
)

func (q *Queue) Cleanup(predicate func(student.Student) bool) []student.Student {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	kept := q.Users[:0]
	removed := make([]student.Student, 0)

	for _, s := range q.Users {
		if predicate(s) {
			removed = append(removed, s)
			continue
		}
		kept = append(kept, s)
	}

	q.Users = kept
	return removed
}

func RemoveIfIdleTooLong(maxIdle time.Duration) func(student.Student) bool {
	return func(s student.Student) bool {
		return time.Since(s.TimeInQueue) >= maxIdle
	}
}

func StartQueueCleanup(
	ctx context.Context,
	bot *telebot.Bot,
	q *Queue,
	interval time.Duration,
	maxIdle time.Duration,
) {
	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return

			case <-ticker.C:
				removed := q.Cleanup(func(s student.Student) bool {
					return time.Since(s.TimeInQueue) >= maxIdle
				})

				for _, s := range removed {
					recipient := &telebot.User{ID: s.ID}

					text := fmt.Sprintf(
						"%s, вы были удалены из очереди за бездействие (>%s). Нажмите «Встать в очередь», чтобы добавиться снова.",
						s.Name,
						maxIdle.Round(time.Minute),
					)

					if _, err := bot.Send(recipient, text); err != nil {
						log.Printf("notify failed for user %d: %v\n", s.ID, err)
					}
				}
			}
		}
	}()
}
