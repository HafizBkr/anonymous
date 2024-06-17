package posts

import (
    "anonymous/models"
    "fmt"
    "github.com/jmoiron/sqlx"
    "time"
    "database/sql"
)
type PostRepo interface {
    CreatePost(payload *PostPayload) (*models.Post, error)
    GetPost(id string) (*models.Post, error)
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
        INSERT INTO posts (id, user_id, content_type, content, description, created_at, likes_count, comments_count)
        VALUES (uuid_generate_v4(), $1, $2, $3, $4, $5, 0, 0)
        RETURNING id
    `
    err := r.db.QueryRow(query, post.UserID, post.ContentType, post.Content, post.Description, post.CreatedAt).Scan(&post.ID)
    if err != nil {
        return nil, fmt.Errorf("error creating post: %w", err)
    }
    return post, nil
}

func (r *postRepo) GetPost(id string) (*models.Post, error) {
    var post models.Post
    query := "SELECT * FROM posts WHERE id = $1"
    err := r.db.Get(&post, query, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, sql.ErrNoRows
        }
        return nil, fmt.Errorf("error getting post: %w", err)
    }
    return &post, nil
}
