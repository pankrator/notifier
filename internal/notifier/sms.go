package notifier

import (
	"context"

	"github.com/pankrator/notifier/internal/config"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type SMSNotifier struct {
	client     *twilio.RestClient
	fromNumber string
}

func NewSMSNotifier(c config.SMSConfig) *SMSNotifier {
	return &SMSNotifier{
		client: twilio.NewRestClientWithParams(
			twilio.ClientParams{
				Username: c.SID,
				Password: c.Secret,
			},
		),
		fromNumber: c.Sender,
	}
}

func (s *SMSNotifier) Send(ctx context.Context, notification Notification) error {
	msg := &twilioApi.CreateMessageParams{}

	msg.SetTo(notification.Recipient)
	msg.SetFrom(s.fromNumber)
	msg.SetBody(notification.Message)

	if _, err := s.client.Api.CreateMessage(msg); err != nil {
		return err
	}

	return nil
}
