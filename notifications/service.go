package notifications

import (
    "fmt"
    "github.com/jmoiron/sqlx"
    "anonymous/types" // Remplacez par le chemin correct pour votre projet
)

type notificationService struct {
    db           *sqlx.DB
    jwtProvider  types.JWTProvider
}

func NewNotificationService(db *sqlx.DB, jwtProvider types.JWTProvider) NotificationService {
    return &notificationService{
        db:          db,
        jwtProvider: jwtProvider,
    }
}

func (s *notificationService) CreateNotification(notification Notification) error {
    _, err := s.db.NamedExec(
        `INSERT INTO notifications (id, user_id, actor_id, action_type, action_id, content, is_read, created_at) VALUES (:id, :user_id, :actor_id, :action_type, :action_id, :content, :is_read, :created_at)`,
        &notification,
    )
    if err != nil {
        return fmt.Errorf("error while creating notification: %w", err)
    }
    return nil
}

func (s *notificationService) DecodeToken(token string) (map[string]interface{}, error) {
    claims, err := s.jwtProvider.Decode(token)
    if err != nil {
        return nil, fmt.Errorf("error decoding token: %w", err)
    }
    return claims, nil
}
