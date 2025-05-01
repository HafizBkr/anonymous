package models

import "time"

type Post struct {
    ID            string    `json:"id" db:"id"`
    UserID        string    `json:"user_id" db:"user_id"`
    Username      string    `json:"username" db:"username"`
    ContentType   string    `json:"content_type" db:"content_type"`
    Content       string    `json:"content" db:"content"`
    Description   string    `json:"description" db:"description"`
    CreatedAt     time.Time `json:"created_at" db:"created_at"`
    LikesCount    int       `json:"likes_count" db:"likes_count"`         // Nombre total de likes
    CommentsCount int       `json:"comments_count" db:"comments_count"`   // Nombre total de commentaires
    LikedByUser   bool      `json:"liked_by_user" db:"liked_by_user"` 
}
