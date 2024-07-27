package models

import "time"

type UserLike struct {
    ID      string    `json:"id" db:"id"`
    UserID  string    `json:"user_id" db:"user_id"`
    LikedBy string    `json:"liked_by" db:"liked_by"`
    LikedAt time.Time `json:"liked_at" db:"liked_at"`
}
