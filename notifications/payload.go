package notifications

type RegisterTokenRequest struct {
    UserID string `json:"user_id"`
    Token  string `json:"token"`
}

type SendNotificationRequest struct {
    Token string `json:"token"`
    Title string `json:"title"`
    Body  string `json:"body"`
}