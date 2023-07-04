package entity

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pankrator/notifier/internal/notifier"
)

type NotificationType string

const (
	EmailNotificationType NotificationType = "EMAIL"
	SMSNotificationType   NotificationType = "SMS"
	SlackNotificationType NotificationType = "SLACK"
)

func (nt NotificationType) ToNotifierType() notifier.NotificationType {
	return notifier.NotificationType(strings.ToLower(string(nt)))
}

type Notification struct {
	ID        uuid.UUID        `db:"id"`
	Type      NotificationType `db:"type"`
	Message   string           `db:"message"`
	Recepient string           `db:"recepient"`
	Metadata  json.RawMessage  `db:"metadata"`
	CreatedAt time.Time        `db:"created_at"`
}
