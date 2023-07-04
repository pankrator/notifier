package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pankrator/notifier/internal/config"
	"github.com/pankrator/notifier/internal/db"
	"github.com/pankrator/notifier/internal/signal"
)

func main() {
	if os.Getenv("LOADED") != "true" {
		if err := godotenv.Load(".env"); err != nil {
			panic(err)
		}
	}

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

	_, err = conn.ExecContext(ctx, `
DROP TYPE IF EXISTS notification_type;
	`)

	if err != nil {
		panic(err)
	}

	_, err = conn.ExecContext(ctx, `
	CREATE TYPE notification_type AS ENUM ('SMS', 'EMAIL', 'SLACK');
		`)

	if err != nil {
		panic(err)
	}

	_, err = conn.ExecContext(ctx, `

CREATE TABLE IF NOT EXISTS notifications(
	id uuid PRIMARY KEY NOT NULL,
	type notification_type NOT NULL,
	message text NOT NULL,
	recepient varchar(500) NOT NULL,
	metadata jsonb NULL,
	created_at timestamp without time zone NOT NULL
);
	`)

	if err != nil {
		panic(err)
	}

	log.Print("Finished migrations")
}
