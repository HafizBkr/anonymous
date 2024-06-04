package posts

import (
    "anonymous/models"
    "fmt"
    "time"

    "github.com/jmoiron/sqlx"
)

type PostRepo interface {
    CreatePost(payload *PostPayload) (*models.Post, error)
}

type postRepo struct {
    db *sqlx.DB
}

func NewPostRepo(db *sqlx.DB) PostRepo {
    return &postRepo{
        db: db,
    }
}

func (r *postRepo) CreatePost(payload *PostPayload) (*models.Post, error) {
    post := &models.Post{
        UserID:      payload.UserID,
        ContentType: payload.ContentType,
        Content:     payload.Content,
        Description: payload.Description,
        CreatedAt:   time.Now(),
    }
    query := `
        INSERT INTO posts (user_id, content_type, content, description, created_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `
    err := r.db.QueryRow(query, post.UserID, post.ContentType, post.Content, post.Description, post.CreatedAt).Scan(&post.ID)
    if err != nil {
        return nil, fmt.Errorf("error creating post: %w", err)
    }
    return post, nil
}
