package points

import (
    "encoding/json"
    "net/http"
    "github.com/go-chi/chi/v5"
    "anonymous/utils"
    "anonymous/types"
)

type PointsHandler struct {
    service PointsService
    logger  types.Logger
}

func NewPointsHandler(service PointsService, logger types.Logger) *PointsHandler {
    return &PointsHandler{service: service, logger: logger}
}

func (h *PointsHandler) HandleLikeUserProfile(w http.ResponseWriter, r *http.Request) {
    // Extraire le token d'autorisation
    token := r.Header.Get("Authorization")
    if token == "" {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Extraire les données du payload
    var payload struct {
        UserID string `json:"user_id"`
    }
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Décoder le token pour obtenir l'ID de l'utilisateur
    claims, err := h.service.DecodeToken(token)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Récupérer l'ID de l'utilisateur à partir des claims
    userID, ok := claims["id"].(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Appeler le service pour liker le profil utilisateur
    err = h.service.LikeUserProfile(payload.UserID, userID)
    if err != nil {
        if err.Error() == "user has already liked this profile" {
            http.Error(w, "User has already liked this profile", http.StatusConflict)
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }
     
    jsonResponse(w, http.StatusOK, nil)
}

func (h *PointsHandler) HandleGetUserProfileLikes(w http.ResponseWriter, r *http.Request) {
    userID := chi.URLParam(r, "userID")

    if userID == "" {
        h.logger.Error("UserID is empty")
        http.Error(w, "UserID is required", http.StatusBadRequest)
        return
    }

    count, err := h.service.GetUserProfileLikes(userID)
    if err != nil {
        h.logger.Error("Error getting user profile likes: " + err.Error())
        utils.WriteError(w, err)
        return
    }

    h.logger.Info("User profile likes retrieved successfully")
    utils.WriteData(w, http.StatusOK, map[string]int{"likes": count})
}

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if data != nil {
        if err := json.NewEncoder(w).Encode(data); err != nil {
            http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        }
    }
}
