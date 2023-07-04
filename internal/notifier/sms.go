package notifier

import (
	"context"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type SMSNotifier struct {
	client     *twilio.RestClient
	fromNumber string
}

func (s *SMSNotifier) Send(ctx context.Context, notification Notification) error {
	msg := &twilioApi.CreateMessageParams{}

	msg.SetTo(notification.Recepient)
	msg.SetFrom(s.fromNumber)
	msg.SetBody(notification.Message)

	if _, err := s.client.Api.CreateMessage(msg); err != nil {
		return err
	}

	return nil
}
