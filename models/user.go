package models

import (
	"time"
    "database/sql"
)


type User struct {
    ID                    string    `db:"id" json:"id"`
    Username              string    `db:"username" json:"username"`
    Password              string    `db:"password_hash" json:"-"`
    Email                 string    `db:"email" json:"email"`
    EmailVerified         bool      `db:"email_verified" json:"email_verified"`
    JoinedAt              time.Time `db:"joined_at" json:"joined_at"`
    Active                bool      `db:"active" json:"active"`
    ProfilePicture        string    `db:"profile_picture" json:"profile_picture"`
    EmailVerificationToken string    `db:"email_verification_token" json:"-"`
    PasswordResetToken    sql.NullString `db:"password_reset_token"`
}

type LoggedInUser struct {
	User
	Token string `json:"token"`
}

