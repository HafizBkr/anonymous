package notifications

import (
    "net/http"
    "github.com/go-chi/render"
)

type NotificationHandler struct {
    service NotificationService
}

func NewNotificationHandler(service NotificationService) *NotificationHandler {
    return &NotificationHandler{service: service}
}

func (h *NotificationHandler) RegisterToken(w http.ResponseWriter, r *http.Request) {
    var request RegisterTokenRequest
    if err := render.DecodeJSON(r.Body, &request); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    if err := h.service.SaveToken(request.UserID, request.Token); err != nil {
        http.Error(w, "Failed to save token", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func (h *NotificationHandler) SendNotification(w http.ResponseWriter, r *http.Request) {
    // Logique pour envoyer une notification push
    var request SendNotificationRequest
    if err := render.DecodeJSON(r.Body, &request); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    err := h.service.SendPushNotification(request.Token, request.Title, request.Body)
    if err != nil {
        http.Error(w, "Failed to send notification", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
