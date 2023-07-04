package notifier

import (
	"context"
	"errors"
)

type NotificationType string

const (
	EmailNotificationType NotificationType = "email"
	SlackNotificationType NotificationType = "slack"
	SMSNotificationType   NotificationType = "sms"
)

var (
	ErrNotifierNotFound = errors.New("no such notifier found")
)

type Client interface {
	Send(context.Context, Notification) error
}

type Notification struct {
	Message   string
	Recipient string
}

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

func (n *Notifier) Send(ctx context.Context, notification Notification, notificationType NotificationType) error {
	client, found := n.clients[notificationType]
	if !found {
		return ErrNotifierNotFound
	}

	return client.Send(ctx, notification)
}
