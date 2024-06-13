package posts

import (
    "encoding/json"
    "net/http"
    "github.com/go-chi/chi/v5"
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

        jsonResponse(w, http.StatusCreated, post)
    }
}
 

func GetPostHandler(service PostService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Récupérer l'ID à partir de l'URL
        id := chi.URLParam(r, "id")

        post, err := service.GetPost(id)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }

        jsonResponse(w, http.StatusOK, post)
    }
}

func updatePostHandler(service PostService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id := chi.URLParam(r, "id") // Récupérer l'ID à partir de l'URL

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

        // Validation du payload du post
        if validationErrs := payload.Validate(); len(validationErrs) > 0 {
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(validationErrs)
            return
        }

        // Mise à jour du post avec l'ID extrait de l'URL et le contenu JSON
        post, err := service.UpdatePost(token, id, &payload)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        jsonResponse(w, http.StatusOK, post)
    }
}

func DeletePostHandler(service PostService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Récupérer l'ID à partir de l'URL
        id := chi.URLParam(r, "id")

        err := service.DeletePost(token, id)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}
func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}
