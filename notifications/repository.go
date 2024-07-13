package notifications

import (
    "github.com/jmoiron/sqlx"
)

type FCMRepo struct {
    db *sqlx.DB
}

func NewFCMRepo(db *sqlx.DB) *FCMRepo {
    return &FCMRepo{db: db}
}

func (r *FCMRepo) SaveToken(userID, token string) error {
    _, err := r.db.Exec("INSERT INTO fcm_tokens (user_id, fcm_token) VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE SET fcm_token = EXCLUDED.fcm_token", userID, token)
    return err
}

func (r *FCMRepo) GetToken(userID string) (string, error) {
    var token string
    err := r.db.Get(&token, "SELECT fcm_token FROM fcm_tokens WHERE user_id = $1", userID)
    return token, err
}

