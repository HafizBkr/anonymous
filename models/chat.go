package models

import (
	"time"
)

type Message struct {
	ID           string    `json:"id" db:"id"`
	From         string    `json:"from" db:"from_user_id"`
	To           string    `json:"to" db:"to_user_id"`
	Content      string    `json:"content" db:"content"`
	SentAt       time.Time `json:"sent_at" db:"sent_at"`
	Owner        bool      `json:"owner,omitempty" db:"-"`
	FromUsername string    `json:"from_username,omitempty" db:"from_user_username"`
	ToUsername   string    `json:"to_username,omitempty" db:"to_user_username"`
}
