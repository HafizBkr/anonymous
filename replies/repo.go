package replies
import (
    "anonymous/models"
    "fmt"
    "github.com/jmoiron/sqlx"
    "time"
    "database/sql"
)

type CommentReplyRepo interface {
    CreateCommentReply(payload *CommentReplyPayload) (*models.CommentReply, error)
    GetCommentRepliesByCommentID(commentID string) ([]*models.CommentReply, error)
    GetCommentReply(id string) (*models.CommentReply, error)
    UpdateCommentReply(id string, payload *UpdateCommentReplyPayload, userID string) (*models.CommentReply, error)
    DeleteCommentReply(userID string, id string) error
}

type commentReplyRepo struct {
    db *sqlx.DB
}

func NewCommentReplyRepo(db *sqlx.DB) CommentReplyRepo {
    return &commentReplyRepo{db: db}
}

func (r *commentReplyRepo) CreateCommentReply(payload *CommentReplyPayload) (*models.CommentReply, error) {
    if validationErrors := payload.Validate(); len(validationErrors) > 0 {
        return nil, fmt.Errorf("validation error: %v", validationErrors)
    }

    reply := &models.CommentReply{
        UserID:      payload.UserID,
        CommentID:   payload.CommentID,
        ContentType: payload.ContentType,
        Content:     payload.Content,
        CreatedAt:   time.Now(),
    }

    query := `
        INSERT INTO comment_replies (id, user_id, comment_id, content_type, content, created_at)
        VALUES (uuid_generate_v4(), $1, $2, $3, $4, $5)
        RETURNING id
    `
    err := r.db.QueryRow(query, reply.UserID, reply.CommentID, reply.ContentType, reply.Content, reply.CreatedAt).Scan(&reply.ID)
    if err != nil {
        return nil, fmt.Errorf("error creating comment reply: %w", err)
    }

    return reply, nil
}

func (r *commentReplyRepo) GetCommentRepliesByCommentID(commentID string) ([]*models.CommentReply, error) {
    var replies []*models.CommentReply
    query := `
        SELECT 
            cr.*, 
            u.username 
        FROM 
            comment_replies cr
        JOIN 
            users u ON cr.user_id = u.id
        WHERE 
            cr.comment_id = $1 
        ORDER BY 
            cr.created_at ASC
    `
    err := r.db.Select(&replies, query, commentID)
    if err != nil {
        return nil, fmt.Errorf("error getting comment replies: %w", err)
    }
    return replies, nil
}


func (r *commentReplyRepo) GetCommentReply(id string) (*models.CommentReply, error) {
    var reply models.CommentReply
    query := `
        SELECT 
            cr.*, 
            u.username 
        FROM 
            comment_replies cr
        JOIN 
            users u ON cr.user_id = u.id
        WHERE 
            cr.id = $1
    `
    err := r.db.Get(&reply, query, id)
    if err != nil {
        return nil, fmt.Errorf("error getting comment reply: %w", err)
    }
    return &reply, nil
}


func (r *commentReplyRepo) UpdateCommentReply(id string, payload *UpdateCommentReplyPayload, userID string) (*models.CommentReply, error) {
    if validationErrors := payload.Validate(); len(validationErrors) > 0 {
        return nil, fmt.Errorf("validation error: %v", validationErrors)
    }

    replyUserID, err := r.GetCommentReplyUserID(id)
    if err != nil {
        return nil, fmt.Errorf("error verifying reply ownership: %w", err)
    }

    if replyUserID != userID {
        return nil, fmt.Errorf("unauthorized: you are not the owner of this reply")
    }

    reply := &models.CommentReply{
        ID:      id,
        Content: payload.Content,
    }

    query := `
        UPDATE comment_replies
        SET content = $1
        WHERE id = $2
        RETURNING id, user_id, comment_id, content_type, content, created_at
    `
    err = r.db.QueryRow(query, reply.Content, reply.ID).Scan(
        &reply.ID, &reply.UserID, &reply.CommentID, &reply.ContentType, &reply.Content, &reply.CreatedAt)
    if err != nil {
        return nil, fmt.Errorf("error updating comment reply: %w", err)
    }

    return reply, nil
}

func (r *commentReplyRepo) GetCommentReplyUserID(replyID string) (string, error) {
    var userID string
    query := "SELECT user_id FROM comment_replies WHERE id = $1"
    err := r.db.Get(&userID, query, replyID)
    if err != nil {
        return "", fmt.Errorf("error getting comment reply user ID: %w", err)
    }
    return userID, nil
}

func (r *commentReplyRepo) DeleteCommentReply(userID string, replyID string) error {
    reply, err := r.GetCommentReply(replyID)
    if err != nil {
        return err
    }

    commentUserID, err := r.GetCommentUserID(reply.CommentID)
    if err != nil {
        return err
    }

    if userID != reply.UserID && userID != commentUserID {
        return fmt.Errorf("unauthorized: you do not have permission to delete this reply")
    }

    deleteQuery := "DELETE FROM comment_replies WHERE id = $1"
    _, err = r.db.Exec(deleteQuery, replyID)
    if err != nil {
        return fmt.Errorf("error deleting comment reply: %w", err)
    }
    return nil
}

func (r *commentReplyRepo) GetCommentUserID(commentID string) (string, error) {
    var userID string
    query := "SELECT user_id FROM comments WHERE id = $1"
    err := r.db.Get(&userID, query, commentID)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", fmt.Errorf("comment not found")
        }
        return "", fmt.Errorf("error getting comment user ID: %w", err)
    }
    return userID, nil
}
