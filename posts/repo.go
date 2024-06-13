package posts

import (
    "anonymous/models"
    "fmt"
    "github.com/jmoiron/sqlx"
    "time"
)

type PostRepo interface {
    CreatePost(payload *PostPayload) (*models.Post, error)
    GetPost(id string) (*models.Post, error)
    UpdatePost(id string, payload *PostPayload) (*models.Post, error)
    DeletePost(userID string, id string) error
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
        return nil, fmt.Errorf("error getting post: %w", err)
    }
    return &post, nil
}

func (r *postRepo) UpdatePost(id string, payload *PostPayload) (*models.Post, error) {
    post := &models.Post{
        ID:          id,
        UserID:      payload.UserID,
        ContentType: payload.ContentType,
        Content:     payload.Content,
        Description: payload.Description,
    }
    query := `
        UPDATE posts SET content_type = $1, content = $2, description = $3
        WHERE id = $4 AND user_id = $5
        RETURNING id, user_id, content_type, content, description, created_at, likes_count, comments_count
    `
    err := r.db.QueryRow(query, post.ContentType, post.Content, post.Description, post.ID, post.UserID).Scan(
        &post.ID, &post.UserID, &post.ContentType, &post.Content, &post.Description, &post.CreatedAt, &post.LikesCount, &post.CommentsCount)
    if err != nil {
        return nil, fmt.Errorf("error updating post: %w", err)
    }
    return post, nil
}

func (r *postRepo) DeletePost(userID string, id string) error {
    query := "DELETE FROM posts WHERE id = $1 AND user_id = $2"
    _, err := r.db.Exec(query, id, userID)
    if err != nil {
        return fmt.Errorf("error deleting post: %w", err)
    }
    return nil
}
