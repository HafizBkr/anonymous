package chat

import (
	middlewares "anonymous/middleware"
	"anonymous/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

func GetConversationsHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(middlewares.ContextKeyUser).(*models.LoggedInUser)
		if user == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		mr := NewMessageRepository(db)
		conversations, err := mr.GetConversations(user.ID)
		if err != nil {
			log.Printf("Error retrieving conversations: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(conversations)
	}
}
