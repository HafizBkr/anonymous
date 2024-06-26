package comments

import (
    "encoding/json"
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/render"
    "fmt"
    "strings"
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
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")
        commentID := chi.URLParam(r, "commentID")

        var payload UpdateCommentPayload
        if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        if validationErrs := payload.Validate(); len(validationErrs) > 0 {
            http.Error(w, fmt.Sprintf("Validation error: %v", validationErrs), http.StatusBadRequest)
            return
        }

        updatedComment, err := service.UpdateComment(token, commentID, &payload)
        if err != nil {
            http.Error(w, fmt.Sprintf("Error updating comment: %v", err), http.StatusInternalServerError)
            return
        }

        render.JSON(w, r, updatedComment)
    }
}



func DeleteCommentHandler(service CommentService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")
        commentID := chi.URLParam(r, "commentID")

        err := service.DeleteComment(token, commentID)
        if err != nil {
            if err.Error() == "comment not found" {
                http.Error(w, "Comment not found", http.StatusNotFound)
                return
            }
            if err.Error() == "post not found" {
                http.Error(w, "Post not found", http.StatusNotFound)
                return
            }
            if err.Error() == "unauthorized: you do not have permission to delete this comment" {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}


