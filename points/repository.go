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
    HasUserLikedProfile(userID, likedBy string) (bool, error)
}

type pointsRepo struct {
    db *sqlx.DB
}

func NewPointsRepo(db *sqlx.DB) PointsRepo {
    return &pointsRepo{db: db}
}

func (r *pointsRepo) LikeUserProfile(userID, likedBy string) error {
    // Vérifiez si l'utilisateur a déjà liké ce profil
    liked, err := r.HasUserLikedProfile(userID, likedBy)
    if err != nil {
        return fmt.Errorf("error while checking if user has liked the profile: %w", err)
    }
    if liked {
        return fmt.Errorf("user has already liked this profile")
    }

    like := models.UserLike{
        ID:      uuid.NewString(),
        UserID:  userID,
        LikedBy: likedBy,
        LikedAt: time.Now(),
    }

    _, err = r.db.NamedExec(
        `INSERT INTO user_likes (id, user_id, liked_by, liked_at) VALUES (:id, :user_id, :liked_by, :liked_at)`,
        &like,
    )
    if err != nil {
        return fmt.Errorf("error while liking user profile: %w", err)
    }
    return nil
}

func (r *pointsRepo) HasUserLikedProfile(userID, likedBy string) (bool, error) {
    var exists bool
    err := r.db.Get(&exists, `SELECT EXISTS(SELECT 1 FROM user_likes WHERE user_id = $1 AND liked_by = $2)`, userID, likedBy)
    if err != nil {
        return false, fmt.Errorf("error while checking if user has liked the profile: %w", err)
    }
    return exists, nil
}

func (r *pointsRepo) CountUserProfileLikes(userID string) (int, error) {
    var count int
    err := r.db.Get(&count, `SELECT COUNT(*) FROM user_likes WHERE user_id = $1`, userID)
    if err != nil {
        return 0, fmt.Errorf("error while counting user profile likes: %w", err)
    }
    return count, nil
}
