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
    GetAllPosts(offset, limit int) ([]*models.Post, error)
    GetPostsByUser(userID string) ([]*models.Post, error)
    UpdatePost(postID string, payload *PostPayload) (*models.Post, error)
    DeletePost(postID string) error
    LikePost(postID, userID string) error
    UnlikePost(postID, userID string) error
    AddReaction(postID, userID, reactionType string) error
    RemoveReaction(postID, userID string) error
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

func (r *postRepo) GetAllPosts(offset, limit int) ([]*models.Post, error) {
    var posts []*models.Post
    query := `
        SELECT 
            p.*, 
            u.username 
        FROM 
            posts p
        JOIN 
            users u ON p.user_id = u.id
        ORDER BY 
            p.created_at DESC 
        LIMIT 
            $1 OFFSET $2
    `
    err := r.db.Select(&posts, query, limit, offset)
    if err != nil {
        return nil, fmt.Errorf("error getting posts: %w", err)
    }
    return posts, nil
}



func (r *postRepo) GetPostsByUser(userID string) ([]*models.Post, error) {
    var posts []*models.Post
    query := `
        SELECT 
            p.*, 
            u.username 
        FROM 
            posts p
        JOIN 
            users u ON p.user_id = u.id
        WHERE 
            p.user_id = $1 
        ORDER BY 
            p.created_at DESC
    `
    err := r.db.Select(&posts, query, userID)
    if err != nil {
        return nil, fmt.Errorf("error getting posts by user: %w", err)
    }
    return posts, nil
}

func (r *postRepo) UpdatePost(postID string, payload *PostPayload) (*models.Post, error) {
    query := `
        UPDATE posts SET content_type = $1, content = $2, description = $3 WHERE id = $4 RETURNING id, user_id, content_type, content, description, created_at, likes_count, comments_count
    `
    var post models.Post
    err := r.db.QueryRow(query, payload.ContentType, payload.Content, payload.Description, postID).Scan(
        &post.ID, &post.UserID, &post.ContentType, &post.Content, &post.Description, &post.CreatedAt, &post.LikesCount, &post.CommentsCount,
    )
    if err != nil {
        return nil, fmt.Errorf("error updating post: %w", err)
    }
    return &post, nil
}

func (r *postRepo) DeletePost(postID string) error {
    tx, err := r.db.Beginx()
    if err != nil {
        return fmt.Errorf("error starting transaction: %w", err)
    }

    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p)
        } else if err != nil {
            tx.Rollback()
        } else {
            tx.Commit()
        }
    }()

    _, err = tx.Exec("DELETE FROM comments WHERE post_id = $1", postID)
    if err != nil {
        return fmt.Errorf("error deleting comments: %w", err)
    }

    _, err = tx.Exec("DELETE FROM posts WHERE id = $1", postID)
    if err != nil {
        return fmt.Errorf("error deleting post: %w", err)
    }

    return nil
}

func (r *postRepo) LikePost(postID, userID string) error {
    query := `INSERT INTO post_likes (post_id, user_id) VALUES ($1, $2)`
    _, err := r.db.Exec(query, postID, userID)
    return err
}

func (r *postRepo) UnlikePost(postID, userID string) error {
    query := `DELETE FROM post_likes WHERE post_id = $1 AND user_id = $2`
    _, err := r.db.Exec(query, postID, userID)
    return err
}

func (r *postRepo) AddReaction(postID, userID, reactionType string) error {
    query := `INSERT INTO post_reactions (post_id, user_id, reaction_type) VALUES ($1, $2, $3)`
    _, err := r.db.Exec(query, postID, userID, reactionType)
    return err
}

func (r *postRepo) RemoveReaction(postID, userID string) error {
    query := `DELETE FROM post_reactions WHERE post_id = $1 AND user_id = $2`
    _, err := r.db.Exec(query, postID, userID)
    return err
}