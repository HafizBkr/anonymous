package notifications

import (
    "context"
    "firebase.google.com/go/messaging"
    "firebase.google.com/go"
    "log"
)

type NotificationService struct {
    app   *firebase.App
    repo  *FCMRepo 
}

func NewNotificationService(app *firebase.App, repo *FCMRepo) *NotificationService {
    return &NotificationService{app: app, repo: repo}
}

func (s *NotificationService) SaveToken(userID, token string) error {
    return s.repo.SaveToken(userID, token)
}

func (s *NotificationService) GetToken(userID string) (string, error) {
    return s.repo.GetToken(userID)
}

func (s *NotificationService) SendPushNotification(token, title, body string) error {
    ctx := context.Background()
    client, err := s.app.Messaging(ctx)
    if err != nil {
        return err
    }

    message := &messaging.Message{
        Notification: &messaging.Notification{
            Title: title,
            Body:  body,
        },
        Token: token,
    }

    _, err = client.Send(ctx, message)
    if err != nil {
        log.Printf("Error sending push notification: %v\n", err)
    }

    return err
}
