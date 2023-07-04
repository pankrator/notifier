package notifier

import (
	"context"

	"github.com/pankrator/notifier/internal/config"
	"gopkg.in/gomail.v2"
)

type Emailer struct {
	client *gomail.Dialer
	sender string
}

func NewEmailer(c config.EmailerConfig) *Emailer {
	return &Emailer{
		client: gomail.NewDialer(
			c.Host,
			c.Port,
			c.Username,
			c.Password,
		),
		sender: c.Sender,
	}
}

func (e *Emailer) Send(ctx context.Context, notification Notification) error {
	m := gomail.NewMessage()

	m.SetHeader("From", e.sender)
	m.SetHeader("To", notification.Recepient)
	m.SetHeader("Subject", "Alert")
	m.SetBody("text/plain", notification.Message)

	if err := e.client.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
