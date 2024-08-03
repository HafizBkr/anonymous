package posts

import (
	"anonymous/auth"
	"anonymous/commons"
	"anonymous/models"
	"anonymous/validator"
	"fmt"
)

type PostService interface {
	CreatePost(token string, payload *PostPayload) (*models.Post, error)
	GetPost(id string) (*models.Post, error)
    GetAllPosts() ([]*models.Post, error)
    GetPostsByUser(userID string) ([]*models.Post, error)
    UpdatePost(token string, postID string, payload *PostPayload) (*models.Post, error)
    DeletePost(token string, postID string) error
    LikePost(token, postID string) error
    UnlikePost(token, postID string) error
    AddReaction(token, postID, reactionType string) error
    RemoveReaction(token, postID string) error
}

type postService struct {
	repo        PostRepo
	authService auth.AuthService
}

func NewPostService(repo PostRepo, authService auth.AuthService) PostService {
	return &postService{
		repo:        repo,
		authService: authService,
	}
}

func (s *postService) CreatePost(token string, payload *PostPayload) (*models.Post, error) {
	userID, err := s.authService.ValidateToken(token)
	if err != nil {
		return nil, commons.Errors.AuthenticationFailed
	}
	payload.UserID = userID

	if err := validator.ValidateStruct(payload); err != nil {
		return nil, err
	}

	return s.repo.CreatePost(payload)
}

func (s *postService) GetPost(id string) (*models.Post, error) {
	return s.repo.GetPost(id)
}

func (s *postService) GetAllPosts() ([]*models.Post, error) {
    return s.repo.GetAllPosts()
}

func (s *postService) GetPostsByUser(userID string) ([]*models.Post, error) {
    return s.repo.GetPostsByUser(userID)
}

func (s *postService) UpdatePost(token string, postID string, payload *PostPayload) (*models.Post, error) {
    userID, err := s.authService.ValidateToken(token)
    if err != nil {
        return nil, commons.Errors.AuthenticationFailed
    }

    if err := validator.ValidateStruct(payload); err != nil {
        return nil, err
    }

    post, err := s.repo.GetPost(postID)
    if err != nil {
        return nil, err
    }

    if post.UserID != userID {
        return nil, fmt.Errorf("unauthorized")
    }

    return s.repo.UpdatePost(postID, payload)
}

func (s *postService) DeletePost(token string, postID string) error {
    userID, err := s.authService.ValidateToken(token)
    if err != nil {
        return commons.Errors.AuthenticationFailed
    }

    post, err := s.repo.GetPost(postID)
    if err != nil {
        return err
    }

    if post.UserID != userID {
        return fmt.Errorf("unauthorized")
    }

    return s.repo.DeletePost(postID)
}

func (s *postService) LikePost(token, postID string) error {
    userID, err := s.authService.ValidateToken(token)
    if err != nil {
        return commons.Errors.AuthenticationFailed
    }

    return s.repo.LikePost(postID, userID)
}

func (s *postService) UnlikePost(token, postID string) error {
    userID, err := s.authService.ValidateToken(token)
    if err != nil {
        return commons.Errors.AuthenticationFailed
    }

    return s.repo.UnlikePost(postID, userID)
}

func (s *postService) AddReaction(token, postID, reactionType string) error {
    userID, err := s.authService.ValidateToken(token)
    if err != nil {
        return commons.Errors.AuthenticationFailed
    }

    return s.repo.AddReaction(postID, userID, reactionType)
}

func (s *postService) RemoveReaction(token, postID string) error {
    userID, err := s.authService.ValidateToken(token)
    if err != nil {
        return commons.Errors.AuthenticationFailed
    }

    return s.repo.RemoveReaction(postID, userID)
}