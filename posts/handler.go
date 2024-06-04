package posts

import (
    "encoding/json"
    "net/http"
)

func CreatePostHandler(service PostService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        var payload PostPayload
        if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        post, err := service.CreatePost(token, &payload)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Return the created post as JSON response
        jsonResponse(w, http.StatusCreated, post)
    }
}

// Helper function to send JSON response
func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}
