package models

import "time"

type PostLike struct {
    ID        string    `db:"id" json:"id"`
    PostID    string    `db:"post_id" json:"post_id"`
    UserID    string    `db:"user_id" json:"user_id"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type PostReaction struct {
    ID           string    `db:"id" json:"id"`
    PostID       string    `db:"post_id" json:"post_id"`
    UserID       string    `db:"user_id" json:"user_id"`
    ReactionType string    `db:"reaction_type" json:"reaction_type"`
    CreatedAt    time.Time `db:"created_at" json:"created_at"`
}
