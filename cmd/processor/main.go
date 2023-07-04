package main

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
	"github.com/pankrator/notifier/internal/config"
	"github.com/pankrator/notifier/internal/db"
	"github.com/pankrator/notifier/internal/entity"
	"github.com/pankrator/notifier/internal/notifier"
	"github.com/pankrator/notifier/internal/processor"
	"github.com/pankrator/notifier/internal/signal"
	"github.com/pankrator/notifier/internal/storage"
)

func main() {
	if os.Getenv("LOADED") != "true" {
		if err := godotenv.Load(".env"); err != nil {
			panic(err)
		}
	}

	notificationTypeArg := flag.String("type", "slack", "notification type to process")

	flag.Parse()

	notificationType := parseType(*notificationTypeArg)

	c, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	ctx, cancel := signal.HandleSignal()
	defer cancel()

	conn, err := db.NewDBConn(ctx, c.DB)
	if err != nil {
		panic(err)
	}

	notificationRepo := storage.NewNotificationRepository(conn)

	not := notifier.NewNotifier()
	not.AddNotifierClient(notifier.EmailNotificationType, notifier.NewEmailer(c.EmailerConfig))
	not.AddNotifierClient(notifier.SlackNotificationType, notifier.NewSlacker(c.SlackerConfig))
	not.AddNotifierClient(notifier.SMSNotificationType, notifier.NewSMSNotifier(c.SMSConfig))

	notificationProcessor := processor.NewProcessor(
		conn,
		notificationRepo,
		not,
		c.Processor,
	)

	if err := notificationProcessor.Start(ctx, notificationType); err != nil {
		panic(err)
	}
}

func parseType(arg string) entity.NotificationType {
	switch arg {
	case "email":
		return entity.EmailNotificationType
	case "sms":
		return entity.SMSNotificationType
	case "slack":
		return entity.SlackNotificationType
	}

	return entity.SlackNotificationType
}
