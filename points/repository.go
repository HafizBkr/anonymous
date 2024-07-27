package points

import (
    "fmt"
    "time"
    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
    "anonymous/models"
)

type PointsRepo interface {
    LikeUserProfile(userID, likedBy string) error
    CountUserProfileLikes(userID string) (int, error)
}

type pointsRepo struct {
    db *sqlx.DB
}

func NewPointsRepo(db *sqlx.DB) PointsRepo {
    return &pointsRepo{db: db}
}

func (r *pointsRepo) LikeUserProfile(userID, likedBy string) error {
    like := models.UserLike{
        ID:      uuid.NewString(),
        UserID:  userID,
        LikedBy: likedBy,
        LikedAt: time.Now(),
    }

    _, err := r.db.NamedExec(
        `INSERT INTO user_likes (id, user_id, liked_by, liked_at) VALUES (:id, :user_id, :liked_by, :liked_at)`,
        &like,
    )
    if err != nil {
        return fmt.Errorf("error while liking user profile: %w", err)
    }
    return nil
}

func (r *pointsRepo) CountUserProfileLikes(userID string) (int, error) {
    var count int

    // Vérifiez que userID est correctement formaté ici si nécessaire
    err := r.db.Get(&count, `SELECT COUNT(*) FROM user_likes WHERE user_id = $1`, userID)
    if err != nil {
        return 0, fmt.Errorf("error while counting user profile likes: %w", err)
    }
    return count, nil
}
