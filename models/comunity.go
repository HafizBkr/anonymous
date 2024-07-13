package models

import "time"

type Community struct {
    ID          string    `db:"id" json:"id"`
    Name        string    `db:"name" json:"name"`
    Description string    `db:"description" json:"description"`
    CreatorID   string    `db:"creator_id" json:"creator_id"`
    CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type CommunityMember struct {
    UserID      string    `db:"user_id" json:"user_id"`
    CommunityID string    `db:"community_id" json:"community_id"`
    JoinedAt    time.Time `db:"joined_at" json:"joined_at"`
}
