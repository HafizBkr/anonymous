package fcm

import (
    "context"
    "firebase.google.com/go"
    "firebase.google.com/go/messaging"
    "log"
)

type FCMService struct {
    app *firebase.App
}

func NewFCMService(app *firebase.App) *FCMService {
    return &FCMService{app: app}
}

func (s *FCMService) SendPushNotification(token, title, body string) error {
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
