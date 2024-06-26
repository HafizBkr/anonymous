package comments

import (
    "anonymous/models"
    "fmt"
    "github.com/jmoiron/sqlx"
    "time"
    "database/sql"
)

type CommentRepo interface {
    CreateComment(payload *CommentPayload) (*models.Comment, error)
    GetComment(id string) (*models.Comment, error)
    UpdateComment(id string, payload *UpdateCommentPayload, userID string) (*models.Comment, error) 
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

func (r *commentRepo) UpdateComment(id string, payload *UpdateCommentPayload, userID string) (*models.Comment, error) {
    if validationErrors := payload.Validate(); len(validationErrors) > 0 {
        return nil, fmt.Errorf("validation error: %v", validationErrors)
    }
    commentUserID, err := r.GetCommentUserID(id)
    if err != nil {
        return nil, fmt.Errorf("error verifying comment ownership: %w", err)
    }

    if commentUserID != userID {
        return nil, fmt.Errorf("unauthorized: you are not the owner of this comment")
    }
    comment := &models.Comment{
        ID:      id,
        Content: payload.Content,
    }
    query := `
        UPDATE comments
        SET content = $1
        WHERE id = $2
        RETURNING id, user_id, post_id, content_type, content, created_at
    `
    err = r.db.QueryRow(query, comment.Content, comment.ID).Scan(
        &comment.ID, &comment.UserID, &comment.PostID, &comment.ContentType, &comment.Content, &comment.CreatedAt)
    if err != nil {
        return nil, fmt.Errorf("error updating comment: %w", err)
    }

    return comment, nil
}


func (r *commentRepo) GetCommentUserID(commentID string) (string, error) {
    var userID string
    query := "SELECT user_id FROM comments WHERE id = $1"
    err := r.db.Get(&userID, query, commentID)
    if err != nil {
        return "", fmt.Errorf("error getting comment user ID: %w", err)
    }
    return userID, nil
}

func (r *commentRepo) DeleteComment(userID string, commentID string) error {
    comment, err := r.GetCommentDetails(commentID)
    if err != nil {
        return err
    }
    postUserID, err := r.GetPostUserID(comment.PostID)
    if err != nil {
        return err
    }
    if userID != comment.UserID && userID != postUserID {
        return fmt.Errorf("unauthorized: you do not have permission to delete this comment")
    }
    deleteQuery := "DELETE FROM comments WHERE id = $1"
    _, err = r.db.Exec(deleteQuery, commentID)
    if err != nil {
        return fmt.Errorf("error deleting comment: %w", err)
    }
    return nil
}


func (r *commentRepo) GetPostUserID(postID string) (string, error) {
    var userID string
    query := "SELECT user_id FROM posts WHERE id = $1"
    err := r.db.Get(&userID, query, postID)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", fmt.Errorf("post not found")
        }
        return "", fmt.Errorf("error getting post user ID: %w", err)
    }
    return userID, nil
}


func (r *commentRepo) GetCommentDetails(commentID string) (*models.Comment, error) {
    var comment models.Comment
    query := "SELECT id, user_id, post_id FROM comments WHERE id = $1"
    err := r.db.Get(&comment, query, commentID)
    if err != nil {
        return nil, fmt.Errorf("error getting comment details: %w", err)
    }
    return &comment, nil
}

