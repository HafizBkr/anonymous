package models

import "time"

type Comment struct {
    ID          string    `json:"id"`
    UserID      string    `json:"user_id"` 
    PostID      string    `json:"post_id"` 
    ContentType string    `json:"content_type"` // Type de contenu du commentaire (texte ou vocal)
    Content     string    `json:"content"` 
    CreatedAt   time.Time `json:"created_at"`
}
