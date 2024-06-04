package posts

import (
    "anonymous/auth"
    "anonymous/commons"
    "anonymous/models"
    "anonymous/validator"
)

type PostService interface {
    CreatePost(token string, payload *PostPayload) (*models.Post, error)
}

type postService struct {
    repo       PostRepo
    authService auth.AuthService
}

func NewPostService(repo PostRepo, authService auth.AuthService) PostService {
    return &postService{
        repo:       repo,
        authService: authService,
    }
}

func (s *postService) CreatePost(token string, payload *PostPayload) (*models.Post, error) {
    userID, err := s.authService.ValidateToken(token)
    if err != nil {
        return nil, commons.Errors.AuthenticationFailed
    }
    payload.UserID = userID // Ensure the UserID in payload matches the authenticated user

    err = validator.ValidateStruct(payload)
    if err != nil {
        return nil, err
    }

    return s.repo.CreatePost(payload)
}
