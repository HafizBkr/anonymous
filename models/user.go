package models

import (
	"anonymous/encryption"
	"time"
    "database/sql"
)

type User struct {
	ID                    string    `db:"id" json:"id"`
	EncryptedEmail        string    `db:"email" json:"-"`
	Email                 string    `db:"-" json:"email"`
	Username              string    `db:"username" json:"username"`
	Password              string    `db:"password_hash" json:"-"`
	JoinedAt              time.Time `db:"joined_at" json:"joined_at"`
	Active                bool      `db:"active" json:"active"`
	ProfilePicture        string    `db:"profile_picture" json:"profile_picture"`
	EmailVerified         bool      `db:"email_verified" json:"email_verified"`
	EmailVerificationToken string   `db:"email_verification_token" json:"-"`
	PasswordResetToken    sql.NullString `db:"password_reset_token"`
}

type LoggedInUser struct {
	ID                    string    `json:"id"`
	EncryptedEmail        string    `json:"-"`
	Email                 string    `json:"email"`
	Username              string    `json:"username"`
	Password              string    `json:"-"`
	JoinedAt              time.Time `json:"joined_at"`
	Active                bool      `json:"active"`
	ProfilePicture        string    `json:"profile_picture"`
	EmailVerified         bool      `json:"email_verified"`
	EmailVerificationToken string    `json:"-"`
}

// DecryptEmail déchiffre l'email de l'utilisateur
func (u *User) DecryptEmail() error {
	ee, err := encryption.NewEmailEncryption()
	if err != nil {
		return err
	}
	
	if u.EncryptedEmail != "" {
		email, err := ee.DecryptEmail(u.EncryptedEmail)
		if err != nil {
			return err
		}
		u.Email = email
	}
	
	return nil
}

// EncryptEmail chiffre l'email de l'utilisateur
func (u *User) EncryptEmail() error {
	ee, err := encryption.NewEmailEncryption()
	if err != nil {
		return err
	}
	
	if u.Email != "" && !ee.IsEncrypted(u.Email) {
		encryptedEmail, err := ee.EncryptEmail(u.Email)
		if err != nil {
			return err
		}
		u.EncryptedEmail = encryptedEmail
	}
	
	return nil
}

// DecryptEmail déchiffre l'email pour LoggedInUser
func (u *LoggedInUser) DecryptEmail() error {
	ee, err := encryption.NewEmailEncryption()
	if err != nil {
		return err
	}
	
	if u.EncryptedEmail != "" {
		email, err := ee.DecryptEmail(u.EncryptedEmail)
		if err != nil {
			return err
		}
		u.Email = email
	}
	
	return nil
}