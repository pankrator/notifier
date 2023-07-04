package http

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pankrator/notifier/internal/entity"
)

type NotificationRepository interface {
	InsertOne(ctx context.Context, not *entity.Notification) error
}

type Handler struct {
	notificationRepository NotificationRepository
}

func (h *Handler) HandleSMSNotification(rw http.ResponseWriter, req *http.Request) {
	h.process(rw, req, entity.SMSNotificationType)
}

func (h *Handler) HandleSlackNotification(rw http.ResponseWriter, req *http.Request) {
	h.process(rw, req, entity.SlackNotificationType)
}

func (h *Handler) HandleEmailNotification(rw http.ResponseWriter, req *http.Request) {
	h.process(rw, req, entity.EmailNotificationType)
}

func (h *Handler) process(rw http.ResponseWriter, req *http.Request, notificationType entity.NotificationType) {
	if req.Method != http.MethodPost {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	type request struct {
		Message   string `json:"message"`
		Recipient string `json:"recipient"`
	}

	var r request

	bytes, err := io.ReadAll(req.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not read body: %s", err)

		return
	}

	if err := json.Unmarshal(bytes, &r); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not unmarshal request: %s", err)

		return
	}

	if err := h.notificationRepository.InsertOne(req.Context(), &entity.Notification{
		ID:        uuid.New(),
		Type:      notificationType,
		Message:   r.Message,
		Recipient: r.Recipient,
		Metadata:  []byte("{}"),
		CreatedAt: time.Now().UTC(),
	}); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not insert notification: %s", err)

		return
	}

	// nolint:errcheck
	rw.Write([]byte("OK"))
}

func NewMuxer(notificationRepository NotificationRepository) *http.ServeMux {
	h := &Handler{
		notificationRepository: notificationRepository,
	}

	muxer := http.NewServeMux()

	muxer.HandleFunc("/sms", recoverHandler(h.HandleSMSNotification))
	muxer.HandleFunc("/email", recoverHandler(h.HandleEmailNotification))
	muxer.HandleFunc("/slack", recoverHandler(h.HandleSlackNotification))

	return muxer
}
