package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/pankrator/notifier/internal/config"
)

type Server struct {
	server *http.Server
	config *config.Server
}

func NewServer(c *config.Server, handler http.Handler) *Server {
	return &Server{
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", c.Host, c.Port),
			Handler: handler,
		},
		config: c,
	}
}

func (s *Server) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		log.Printf("Server listening on address %s:%d", s.config.Host, s.config.Port)

		if err := s.server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}

		log.Printf("Server stopped")
	}()

	go func() {
		<-ctx.Done()

		log.Printf("Shutting down server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.ShutdownTimeout))
		defer cancel()

		if err := s.server.Shutdown(shutdownCtx); err != nil {
			log.Fatal(err)
		}
	}()
}
