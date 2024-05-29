package models

import "time"

type User struct {
    ID            string    `json:"id"`
    Email         string    `json:"email"`
    Username      string    `json:"username"`
    Password      string     `json:"password"`
    CreatedAt     time.Time `json:"created_at"`
    Picture       string    `json:"picture"`
    Active        bool      `json:"active"`           // Champ pour l'état actif/inactif de l'utilisateur
    ProfilePicture string    `json:"profile_picture"` // Champ pour l'image de profil de l'utilisateur
}
