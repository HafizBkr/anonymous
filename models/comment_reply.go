package models

import "time"

type CommentReply struct {
    ID          string    `json:"id" db:"id"`
    UserID      string    `json:"user_id" db:"user_id"` // ID de l'utilisateur qui a créé la réponse au commentaire
    Username    string    `db:"username" json:"username"`
    CommentID   string    `json:"comment_id" db:"comment_id"` // ID du commentaire auquel la réponse est associée
    ContentType string    `json:"content_type" db:"content_type"` // Type de contenu de la réponse au commentaire (texte ou vocal)
    Content     string    `json:"content" db:"content"` // Contenu de la réponse au commentaire (texte ou lien vers le fichier vocal)
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
