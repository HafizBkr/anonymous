package models

import "time"

type ChatMessage struct {
    ID          string    `json:"id"`
    SenderID    string    `json:"sender_id"` // ID de l'utilisateur qui a envoyé le message
    ReceiverID  string    `json:"receiver_id"` // ID de l'utilisateur qui reçoit le message
    ContentType string    `json:"content_type"` // Type de contenu du message (texte ou vocal)
    Content     string    `json:"content"` // Contenu du message (texte ou lien vers le fichier vocal)
    CreatedAt   time.Time `json:"created_at"`
}
