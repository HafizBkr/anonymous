package models

import (
	"time"
)

type Conversation struct {
	UserID             string    `json:"user_id" db:"user_id"`
	Username           string    `json:"username" db:"username"`
	ProfilePicture     string    `json:"profile_picture" db:"profile_picture"`
	LastMessageID      string    `json:"last_message_id" db:"last_message_id"`
	LastMessageContent string    `json:"last_message_content" db:"last_message_content"`
	LastMessageSentAt  time.Time `json:"last_message_sent_at" db:"last_message_sent_at"`
}
