package models

import (
	"time"
)

type User struct {
    ID             string    `json:"id" db:"id"`
    Email          string    `json:"email" db:"email"`
    Username       string    `json:"username" db:"username"`
    Password       string       `json:"password_hash" db:"password_hash"`
    JoinedAt       time.Time `json:"joined_at" db:"joined_at"`
    Active         bool      `json:"active" db:"active"`
    ProfilePicture string    `json:"profile_picture" db:"profile_picture"`
    EmailVerified  bool      `json:"email_verified" db:"email_verified"`
    EmailVerificationToken string `json:"email_verification_token" db:"email_verification_token"`
}


type LoggedInUser struct {
	User
	Token string `json:"token"`
}
