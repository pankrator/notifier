package main

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/pankrator/notifier/internal/config"
	"github.com/pankrator/notifier/internal/db"
	"github.com/pankrator/notifier/internal/http"
	"github.com/pankrator/notifier/internal/server"
	"github.com/pankrator/notifier/internal/signal"
	"github.com/pankrator/notifier/internal/storage"
)

func main() {
	if os.Getenv("LOADED") != "true" {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal(err)
		}
	}

	c, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.HandleSignal()
	defer cancel()

	conn, err := db.NewDBConn(ctx, c.DB)
	if err != nil {
		log.Fatal(err)
	}

	notificationRepo := storage.NewNotificationRepository(conn)

	muxer := http.NewMuxer(notificationRepo)

	s := server.NewServer(&c.Server, muxer)

	wg := &sync.WaitGroup{}

	s.Start(ctx, wg)

	wg.Wait()
}
