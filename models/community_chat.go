package models

import "time"

type CommunityChat struct {
    ID          string    `db:"id" json:"id"`
    CommunityID string    `db:"community_id" json:"community_id"`
    UserID      string    `db:"user_id" json:"user_id"`
    Message     string    `db:"message" json:"message"`
    CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
