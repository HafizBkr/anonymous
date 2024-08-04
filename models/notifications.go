package models

import "time"

type Notification struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	ActorID       string    `json:"actor_id"`
	ActionType    string    `json:"action_type"`
	ActionID      string    `json:"action_id"`
	Content       string    `json:"content"`
	IsRead        bool      `json:"is_read"`
	CreatedAt     time.Time `json:"created_at"`
}
