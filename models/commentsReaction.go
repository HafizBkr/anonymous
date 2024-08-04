package models

import "time"

type CommentReaction struct {
    ID           string    `db:"id" json:"id"`
    CommentID    string    `db:"comment_id" json:"comment_id"`
    UserID       string    `db:"user_id" json:"user_id"`
    ReactionType string    `db:"reaction_type" json:"reaction_type"`
    CreatedAt    time.Time `db:"created_at" json:"created_at"`
}
