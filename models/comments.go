package models

import "time"

type Comment struct {
    ID          string    `db:"id" json:"id"`
    UserID      string    `db:"user_id" json:"user_id"`
    PostID      string    `db:"post_id" json:"post_id"`
    ContentType string    `db:"content_type" json:"content_type"`
    Content     string    `db:"content" json:"content"`
    CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
