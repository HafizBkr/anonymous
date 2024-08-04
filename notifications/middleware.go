package notifications

import (
    "encoding/json"
    "fmt"
    "time"
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"
)

func NotificationObserver(notificationService NotificationService) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            next.ServeHTTP(w, r)

            // Extraire le token d'autorisation
            token := r.Header.Get("Authorization")
            if token == "" {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }

            // Décoder le token pour obtenir les informations de l'utilisateur
            claims, err := notificationService.DecodeToken(token)
            if err != nil {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }

            // Extraire l'ID de l'utilisateur effectuant l'action du token
            actorID, ok := claims["id"].(string)
            if !ok {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }

            // Extraire le user_id du corps de la requête pour l'utilisateur dont le profil est liké
            var payload struct {
                UserID string `json:"user_id"`
            }
            if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
                http.Error(w, "Invalid request payload", http.StatusBadRequest)
                return
            }

            userID := payload.UserID
            actionType := r.Header.Get("X-Action-Type")
            actionID := chi.URLParam(r, "id")
            content := r.Header.Get("X-Content")

            if actionType != "" && userID != "" && actorID != "" {
                notification := Notification{
                    ID:         uuid.New().String(),
                    UserID:     userID,
                    ActorID:    actorID,
                    ActionType: actionType,
                    ActionID:   actionID,
                    Content:    content,
                    CreatedAt:  time.Now(),
                }
                if err := notificationService.CreateNotification(notification); err != nil {
                    fmt.Printf("Error creating notification: %v\n", err)
                }
            }
        })
    }
}
