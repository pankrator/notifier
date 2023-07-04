package notifier

import "context"

type NotificationType string

const (
	EmailNotificationType NotificationType = "email"
	SlackNotificationType NotificationType = "slack"
	SMSNotificationType   NotificationType = "sms"
)

type Notifier struct {
	clients map[NotificationType]Client
}

func NewNotifier() *Notifier {
	return &Notifier{
		clients: make(map[NotificationType]Client),
	}
}

func (n *Notifier) AddNotifierClient(notificationType NotificationType, c Client) {
	n.clients[notificationType] = c
}

type Client interface {
	Send(context.Context, Notification) error
}

type Notification struct {
	Message   string
	Recepient string
}

type EmailNotification struct {
	Message   string
	Recepient string
}

func NewEmailNotification(message, recepient string) *EmailNotification {
	return &EmailNotification{
		Message:   message,
		Recepient: recepient,
	}
}

func (n *Notifier) Send(ctx context.Context, notification Notification, notificationType NotificationType) error {
	return n.clients[notificationType].Send(ctx, notification)
}
