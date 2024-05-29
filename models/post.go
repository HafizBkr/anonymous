package models

import "time"

type Post struct {
    ID          string    `json:"id"`
    UserID      string    `json:"user_id"`
    ContentType string    `json:"content_type"`
    Content     string    `json:"content"`
    Description string    `json:"description"` 
    CreatedAt   time.Time `json:"created_at"`
}
