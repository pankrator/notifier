package service

import "github.com/pankrator/notifier/internal/notifier"

type Notifier interface {
	Send(notification notifier.Notification, notificationType notifier.NotificationType) error
}

type SenderService struct {
	notifier Notifier
}

// func (s *SenderService) SendEmail()
