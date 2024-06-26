package replies

import (
    "encoding/json"
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/render"
    "fmt"
    "strings"   
)

func CreateCommentReplyHandler(service CommentReplyService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        commentID := chi.URLParam(r, "commentID")
        if commentID == "" {
            http.Error(w, "Comment ID is required", http.StatusBadRequest)
            return
        }

        var payload CommentReplyPayload
        if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }
        
        payload.CommentID= commentID


        createdReply, err := service.CreateCommentReply(token, &payload)
        if err != nil {
            http.Error(w, fmt.Sprintf("Failed to create comment reply: %v", err), http.StatusInternalServerError)
            return
        }

        render.JSON(w, r, createdReply)
    }
}

func GetCommentRepliesByCommentIDHandler(service CommentReplyService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        commentID := chi.URLParam(r, "commentID")

        replies, err := service.GetCommentRepliesByCommentID(commentID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        render.JSON(w, r, replies)
    }
}

func GetCommentReplyHandler(service CommentReplyService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        replyID := chi.URLParam(r, "replyID")

        reply, err := service.GetCommentReply(replyID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }

        render.JSON(w, r, reply)
    }
}

func UpdateCommentReplyHandler(service CommentReplyService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")
        replyID := chi.URLParam(r, "replyID")

        var payload UpdateCommentReplyPayload
        if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        if validationErrs := payload.Validate(); len(validationErrs) > 0 {
            http.Error(w, fmt.Sprintf("Validation error: %v", validationErrs), http.StatusBadRequest)
            return
        }

        updatedReply, err := service.UpdateCommentReply(token, replyID, &payload)
        if err != nil {
            http.Error(w, fmt.Sprintf("Error updating reply: %v", err), http.StatusInternalServerError)
            return
        }

        render.JSON(w, r, updatedReply)
    }
}

func DeleteCommentReplyHandler(service CommentReplyService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")
        replyID := chi.URLParam(r, "replyID")

        err := service.DeleteCommentReply(token, replyID)
        if err != nil {
            if err.Error() == "reply not found" {
                http.Error(w, "Reply not found", http.StatusNotFound)
                return
            }
            if err.Error() == "comment not found" {
                http.Error(w, "Comment not found", http.StatusNotFound)
                return
            }
            if err.Error() == "unauthorized: you do not have permission to delete this reply" {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}
