package comments

import (
    "encoding/json"
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/render"
)

func CreateCommentHandler(service CommentService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        postID := chi.URLParam(r, "postID")

        var payload CommentPayload
        if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        // Set postID from URL parameter
        payload.PostID = postID

        comment, err := service.CreateComment(token, &payload)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        render.JSON(w, r, comment)
    }
}


func GetCommentsByPostIDHandler(service CommentService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        postID := chi.URLParam(r, "postID")

        comments, err := service.GetCommentsByPostID(postID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        render.JSON(w, r, comments)
    }
}

func GetCommentHandler(service CommentService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        commentID := chi.URLParam(r, "commentID")

        comment, err := service.GetComment(commentID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }

        render.JSON(w, r, comment)
    }
}

func UpdateCommentHandler(service CommentService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        commentID := chi.URLParam(r, "commentID")
        postID := chi.URLParam(r, "postID") // Récupérer l'ID du post depuis l'URL

        var payload CommentPayload
        if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }
        payload.PostID = postID

        // Validate the payload
        if validationErrs := payload.Validate(); len(validationErrs) > 0 {
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(validationErrs)
            return
        }

        // Proceed with updating the comment
        comment, err := service.UpdateComment(token, commentID, &payload)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        render.JSON(w, r, comment)
    }
}



func DeleteCommentHandler(service CommentService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        commentID := chi.URLParam(r, "commentID")
        if err := service.DeleteComment(token, commentID); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}
