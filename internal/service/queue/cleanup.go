package queue

import (
	"context"
	"fmt"
	"log"
	"time"

	"suai-queue/internal/domain"

	"gopkg.in/telebot.v3"
)

func (q *Queue) Cleanup(predicate func(domain.Student) bool) []domain.Student {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	kept := q.Users[:0]
	removed := make([]domain.Student, 0)

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
				removed := q.Cleanup(func(s domain.Student) bool {
					return time.Since(s.TimeInQueue) >= maxIdle
				})

				for _, s := range removed {
					recipient := &telebot.User{ID: s.TgID}

					text := fmt.Sprintf(
						"%s, вы были удалены из очереди за бездействие (>%s). Нажмите «Встать в очередь», чтобы добавиться снова.",
						s.Name,
						maxIdle.Round(time.Minute),
					)

					if _, err := bot.Send(recipient, text); err != nil {
						log.Printf("notify failed for user %d: %v\n", s.TgID, err)
					}
				}
			}
		}
	}()
}
