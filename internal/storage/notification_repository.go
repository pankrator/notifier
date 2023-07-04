package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pankrator/notifier/internal/entity"
)

type NotificationRepository struct {
	queryer Queryer
}

func NewNotificationRepository(queryer Queryer) *NotificationRepository {
	return &NotificationRepository{
		queryer: queryer,
	}
}

func (r *NotificationRepository) InsertOne(ctx context.Context, not *entity.Notification) error {
	res, err := r.queryer.NamedExecContext(ctx, `INSERT INTO notifications
		(id, type, message, recepient, metadata, created_at)
		VALUES
		(:id, :type, :message, :recepient, :metadata, :created_at)`, not)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n < 1 {
		return errors.New("could not insert notification")
	}

	return nil
}

func (r *NotificationRepository) SelectByTypeForUpdate(ctx context.Context, notificationType entity.NotificationType) ([]*entity.Notification, error) {
	result := make([]*entity.Notification, 0)

	if err := r.queryer.SelectContext(ctx, &result, `
	SELECT id, type, message, recepient, metadata, created_at
	FROM notifications WHERE type=$1
	ORDER BY created_at ASC
	FOR UPDATE SKIP LOCKED
	LIMIT 1
	`, notificationType); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *NotificationRepository) DeleteByIDs(ctx context.Context, ids []uuid.UUID) error {
	query, args, err := sqlx.In("DELETE FROM notifications where id IN (?)", ids)
	if err != nil {
		return err
	}

	query = r.queryer.Rebind(query)
	_, err = r.queryer.ExecContext(ctx, query, args...)
	return err
}