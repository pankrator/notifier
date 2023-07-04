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
	Recipient string           `db:"recipient"`
	// Metadata for other types of notifications more data might be needed that is not generic to all types
	// it can be saved in metadata.
	// Other approach would be to have join tables with additional data for certain types of notifications
	// that need the additional data.
	Metadata  json.RawMessage `db:"metadata"`
	CreatedAt time.Time       `db:"created_at"`
}
