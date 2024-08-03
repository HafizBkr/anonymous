package communitychats

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"anonymous/models"
)

type CommunityChatHandler struct {
	Service CommunityChatService
}

func NewCommunityChatHandler(service CommunityChatService) *CommunityChatHandler {
	return &CommunityChatHandler{Service: service}
}

func (h *CommunityChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	communityID := chi.URLParam(r, "communityID")
	var req struct {
		UserID  string `json:"user_id"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	chat := models.CommunityChat{
		ID:          uuid.New().String(),
		CommunityID: communityID,
		UserID:      req.UserID,
		Message:     req.Message,
		CreatedAt:   time.Now(),
	}

	if err := h.Service.CreateMessage(chat); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *CommunityChatHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	communityID := chi.URLParam(r, "communityID")

	chats, err := h.Service.GetMessagesByCommunityID(communityID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chats); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
