package models

import (
	"time"
)

type Story struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Duration  int       `json:"duration" db:"duration"` // Duration in hours
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
}

type StoryContent struct {
	ID        string    `json:"id" db:"id"`
	StoryID   string    `json:"story_id" db:"story_id"`
	Type      string    `json:"type" db:"type"`       // "text" or "image"
	Content   string    `json:"content" db:"content"` // URL or text content
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type StoryLike struct {
	ID        string    `json:"id" db:"id"`
	StoryID   string    `json:"story_id" db:"story_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type StoryView struct {
	ID        string    `json:"id" db:"id"`
	StoryID   string    `json:"story_id" db:"story_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type StoryResponse struct {
	ID        string    `json:"id" db:"id"`
	StoryID   string    `json:"story_id" db:"story_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
