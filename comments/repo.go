package comments

import (
    "anonymous/models"
    "fmt"
    "github.com/jmoiron/sqlx"
    "time"
)

type CommentRepo interface {
    CreateComment(payload *CommentPayload) (*models.Comment, error)
    GetComment(id string) (*models.Comment, error)
    UpdateComment(id string, payload *CommentPayload) (*models.Comment, error)
    DeleteComment(userID string, id string) error
    GetCommentsByPostID(postID string) ([]*models.Comment, error)
}

type commentRepo struct {
    db *sqlx.DB
}

func NewCommentRepo(db *sqlx.DB) CommentRepo {
    return &commentRepo{db: db}
}

func (r *commentRepo) CreateComment(payload *CommentPayload) (*models.Comment, error) {
    if validationErrors := payload.Validate(); len(validationErrors) > 0 {
        return nil, fmt.Errorf("validation error: %v", validationErrors)
    }

    comment := &models.Comment{
        UserID:      payload.UserID,
        PostID:      payload.PostID,
        ContentType: payload.ContentType,
        Content:     payload.Content,
        CreatedAt:   time.Now(),
    }

    query := `
        INSERT INTO comments (id, user_id, post_id, content_type, content, created_at)
        VALUES (uuid_generate_v4(), $1, $2, $3, $4, $5)
        RETURNING id
    `
    err := r.db.QueryRow(query, comment.UserID, comment.PostID, comment.ContentType, comment.Content, comment.CreatedAt).Scan(&comment.ID)
    if err != nil {
        return nil, fmt.Errorf("error creating comment: %w", err)
    }

    return comment, nil
}

func (r *commentRepo) GetCommentsByPostID(postID string) ([]*models.Comment, error) {
    var comments []*models.Comment
    query := "SELECT * FROM comments WHERE post_id = $1 ORDER BY created_at ASC"
    err := r.db.Select(&comments, query, postID)
    if err != nil {
        return nil, fmt.Errorf("error getting comments: %w", err)
    }
    return comments, nil
}

func (r *commentRepo) GetComment(id string) (*models.Comment, error) {
    var comment models.Comment
    query := "SELECT * FROM comments WHERE id = $1"
    err := r.db.Get(&comment, query, id)
    if err != nil {
        return nil, fmt.Errorf("error getting comment: %w", err)
    }
    return &comment, nil
}

func (r *commentRepo) UpdateComment(id string, payload *CommentPayload) (*models.Comment, error) {
    if validationErrors := payload.Validate(); len(validationErrors) > 0 {
        return nil, fmt.Errorf("validation error: %v", validationErrors)
    }

    comment := &models.Comment{
        ID:          id,
        UserID:      payload.UserID,
        PostID:      payload.PostID,
        ContentType: payload.ContentType,
        Content:     payload.Content,
    }

    query := `
        UPDATE comments SET content_type = $1, content = $2
        WHERE id = $3 AND user_id = $4
        RETURNING id, user_id, post_id, content_type, content, created_at
    `
    err := r.db.QueryRow(query, comment.ContentType, comment.Content, comment.ID, comment.UserID).Scan(
        &comment.ID, &comment.UserID, &comment.PostID, &comment.ContentType, &comment.Content, &comment.CreatedAt)
    if err != nil {
        return nil, fmt.Errorf("error updating comment: %w", err)
    }
    return comment, nil
}

func (r *commentRepo) DeleteComment(userID string, id string) error {
    query := "DELETE FROM comments WHERE id = $1 AND user_id = $2"
    _, err := r.db.Exec(query, id, userID)
    if err != nil {
        return fmt.Errorf("error deleting comment: %w", err)
    }
    return nil
}
