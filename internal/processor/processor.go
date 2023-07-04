package processor

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/pankrator/notifier/internal/config"
	"github.com/pankrator/notifier/internal/db"
	"github.com/pankrator/notifier/internal/entity"
	"github.com/pankrator/notifier/internal/notifier"
	"github.com/pkg/errors"
)

type NotificationRepository interface {
	SelectByTypeForUpdate(context.Context, entity.NotificationType) ([]*entity.Notification, error)
	DeleteByIDs(ctx context.Context, ids []uuid.UUID) error
}

type Notifier interface {
	Send(ctx context.Context, notification notifier.Notification, notificationType notifier.NotificationType) error
}

type Processor struct {
	transactioner          db.Transactioner
	notificationRepository NotificationRepository
	notifier               Notifier
	fetchInterval          time.Duration
}

func NewProcessor(
	transactioner db.Transactioner,
	notificationRepository NotificationRepository,
	notifier Notifier,
	c config.Processor,
) *Processor {
	return &Processor{
		transactioner:          transactioner,
		notificationRepository: notificationRepository,
		notifier:               notifier,
		fetchInterval:          c.FetchInterval,
	}
}

func (p *Processor) Start(ctx context.Context, notificationType entity.NotificationType) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		if err := p.transactioner.RunInTx(ctx, func(ctx context.Context) error {
			result, err := p.notificationRepository.SelectByTypeForUpdate(ctx, notificationType)
			if err != nil {
				log.Printf("Could not select notification: %s", err)
			}

			if len(result) < 1 {
				log.Printf("No notifications to process")
				return nil
			}

			ids := make([]uuid.UUID, 0, len(result))
			for _, emailNot := range result {
				if err := p.notifier.Send(ctx, notifier.Notification{
					Message:   emailNot.Message,
					Recipient: emailNot.Recipient,
				}, notificationType.ToNotifierType()); err != nil {
					return errors.Wrap(err, "could not send notification")
				}

				ids = append(ids, emailNot.ID)
			}

			if err := p.notificationRepository.DeleteByIDs(ctx, ids); err != nil {
				return errors.Wrap(err, "could not delete notifications")
			}

			return nil
		}); err != nil {
			log.Printf("Could not process notifications: %s", err)
		}

		time.Sleep(p.fetchInterval)
	}
}
