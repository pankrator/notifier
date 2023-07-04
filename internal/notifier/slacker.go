package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pankrator/notifier/internal/config"
)

type Slacker struct {
	client     *http.Client
	webhookURL string
}

func NewSlacker(c config.SlackerConfig) *Slacker {
	return &Slacker{
		client:     http.DefaultClient,
		webhookURL: c.Webhook,
	}
}

type slackMessage struct {
	Text string `json:"text"`
}

func (s *Slacker) Send(ctx context.Context, notification Notification) error {
	msg := slackMessage{
		Text: notification.Message,
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal Slack message: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, s.webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("could not close body: %s", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	return nil
}
